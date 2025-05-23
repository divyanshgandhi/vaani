import 'package:dio/dio.dart';
import 'package:vaani/core/network/interceptors/auth_interceptor.dart';
import 'package:vaani/domain/repositories/auth_repository.dart';

class ApiClient {
  final Dio _dio;

  ApiClient(this._dio) {
    // Use environment variable or default to local backend
    const baseUrl = String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'http://localhost:8080/v1',
    );

    _dio.options.baseUrl = baseUrl;
    _dio.options.connectTimeout = const Duration(seconds: 30);
    _dio.options.receiveTimeout = const Duration(seconds: 30);
    _dio.options.headers = {'Content-Type': 'application/json'};

    // Add interceptors for logging, authorization, etc.
    _dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
      requestHeader: true,
      responseHeader: false,
    ));
  }

  // Initialize auth interceptor after auth repository is available
  void initAuthInterceptor(AuthRepository authRepository) {
    _dio.interceptors.add(AuthInterceptor(
        authRepository,
        Dio(BaseOptions(
          baseUrl: _dio.options.baseUrl,
          connectTimeout: _dio.options.connectTimeout,
          receiveTimeout: _dio.options.receiveTimeout,
        ))));
  }

  // HTTP methods
  Future<Response> get(String path, {Map<String, dynamic>? queryParameters}) {
    return _dio.get(path, queryParameters: queryParameters);
  }

  Future<Response> post(String path, {Map<String, dynamic>? data}) {
    return _dio.post(path, data: data);
  }

  Future<Response> put(String path, {Map<String, dynamic>? data}) {
    return _dio.put(path, data: data);
  }

  Future<Response> delete(String path) {
    return _dio.delete(path);
  }

  // WebSocket URL for preview functionality
  String get webSocketBaseUrl {
    final baseUrl = _dio.options.baseUrl;
    if (baseUrl.startsWith('https://')) {
      return baseUrl.replaceFirst('https://', 'wss://');
    } else if (baseUrl.startsWith('http://')) {
      return baseUrl.replaceFirst('http://', 'ws://');
    }
    return 'ws://localhost:8090/v1';
  }
}
