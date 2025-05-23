import '../models/voice.dart';

abstract class VoiceRepository {
  /// Fetches available voices
  Future<List<Voice>> getVoices();

  /// Fetches a single voice by ID
  Future<Voice?> getVoice(String id);
}
