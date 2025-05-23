import 'dart:async';
import 'package:dio/dio.dart';
import 'package:vaani/domain/repositories/auth_repository.dart';

class AuthInterceptor extends Interceptor {
  final AuthRepository _authRepository;
  final Dio _dio;
  bool _isRefreshing = false;
  final List<RequestOptions> _pendingRequests = [];

  AuthInterceptor(this._authRepository, this._dio);

  @override
  void onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    // Skip adding token for authentication endpoints
    if (options.path.contains('auth/') && !options.path.contains('refresh')) {
      return handler.next(options);
    }

    final token = await _authRepository.getAuthToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    return handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    // If error is unauthorized and not currently refreshing token
    if (err.response?.statusCode == 401 && !_isRefreshing) {
      // Save the failed request
      _pendingRequests.add(err.requestOptions);

      // Try to refresh token
      try {
        _isRefreshing = true;

        // Call refresh token API - in a real app, you'd implement this method
        final newToken = await _refreshToken();

        if (newToken != null) {
          await _authRepository.saveAuthToken(newToken);

          // Retry all pending requests with new token
          await _retryFailedRequests(newToken);

          // Complete the original request
          final response = await _retryRequest(err.requestOptions, newToken);
          return handler.resolve(response);
        } else {
          // If refresh failed, logout user
          await _authRepository.logout();
          return handler.next(err);
        }
      } catch (e) {
        await _authRepository.logout();
        return handler.next(err);
      } finally {
        _isRefreshing = false;
        _pendingRequests.clear();
      }
    }

    // For other errors, just pass them through
    return handler.next(err);
  }

  // Helper method to refresh token
  Future<String?> _refreshToken() async {
    try {
      // In a real app, implement token refresh logic here
      // For example:
      final refreshResponse = await _dio.post(
        '/auth/refresh',
        options: Options(
          headers: {
            'Authorization': 'Bearer ${await _authRepository.getAuthToken()}',
          },
        ),
      );

      if (refreshResponse.statusCode == 200) {
        return refreshResponse.data['token'];
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  // Retry failed requests with new token
  Future<void> _retryFailedRequests(String newToken) async {
    for (final request in _pendingRequests) {
      try {
        await _retryRequest(request, newToken);
      } catch (e) {
        // Ignore retry errors
      }
    }
  }

  // Retry a single request with the new token
  Future<Response> _retryRequest(
      RequestOptions request, String newToken) async {
    final options = Options(
      method: request.method,
      headers: {
        ...request.headers,
        'Authorization': 'Bearer $newToken',
      },
    );

    return _dio.request(
      request.path,
      data: request.data,
      queryParameters: request.queryParameters,
      options: options,
    );
  }
}
