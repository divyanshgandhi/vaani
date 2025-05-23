import 'package:dio/dio.dart';
import '../../domain/repositories/pricing_repository.dart';

class PricingRepositoryImpl implements PricingRepository {
  final Dio _dio;

  PricingRepositoryImpl(this._dio);

  @override
  Future<double> getPricePerCharacter() async {
    try {
      final response = await _dio.get('/pricing');
      return response.data['pricePerCharacter'] ?? 0.01;
    } catch (e) {
      // Return a default value if API fails
      return 0.01;
    }
  }
}
