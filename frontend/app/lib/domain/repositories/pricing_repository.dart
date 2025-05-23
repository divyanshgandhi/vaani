abstract class PricingRepository {
  /// Fetches the current pricing information
  ///
  /// Returns the price per character in INR
  Future<double> getPricePerCharacter();
}
