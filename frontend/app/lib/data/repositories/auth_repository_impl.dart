import 'package:shared_preferences/shared_preferences.dart';
import 'package:vaani/core/network/api_client.dart';
import 'package:vaani/domain/repositories/auth_repository.dart';

class AuthRepositoryImpl implements AuthRepository {
  final ApiClient _apiClient;
  final SharedPreferences _prefs;
  static const String _tokenKey = 'auth_token';

  AuthRepositoryImpl(this._apiClient, this._prefs);

  @override
  Future<bool> isLoggedIn() async {
    final token = await getAuthToken();
    return token != null && token.isNotEmpty;
  }

  @override
  Future<String?> getAuthToken() async {
    return _prefs.getString(_tokenKey);
  }

  @override
  Future<void> saveAuthToken(String token) async {
    await _prefs.setString(_tokenKey, token);
  }

  @override
  Future<void> clearAuthToken() async {
    await _prefs.remove(_tokenKey);
  }

  @override
  Future<bool> sendOtp(String phoneNumber) async {
    try {
      final response = await _apiClient.post('/auth/send-otp', data: {
        'phoneNumber': phoneNumber,
      });
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }

  @override
  Future<bool> verifyOtp(String phoneNumber, String otp) async {
    try {
      final response = await _apiClient.post('/auth/verify-otp', data: {
        'phoneNumber': phoneNumber,
        'otp': otp,
      });

      if (response.statusCode == 200 && response.data['token'] != null) {
        await saveAuthToken(response.data['token']);
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  @override
  Future<void> logout() async {
    await clearAuthToken();
  }
}
