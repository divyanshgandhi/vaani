import 'package:vaani/domain/entities/language_detection.dart';
import 'package:vaani/domain/entities/tts_job.dart';
import 'package:vaani/domain/entities/voice.dart';

abstract class TtsRepository {
  /// Detects the language of the given text
  Future<LanguageDetection> detectLanguage(String text);

  /// Generates TTS audio from the given text
  Future<TtsJob> generateTts({
    required String text,
    required String voiceId,
    String? language,
    int emotion = 50,
    double speed = 1.0,
    double pitch = 1.0,
    String outputType = 'mp3',
  });

  /// Gets the status of a TTS job
  Future<TtsJob> getJobStatus(String jobId);

  /// Gets the download URL for a completed TTS job
  Future<String> getDownloadUrl(String jobId);

  /// Gets the list of available voices
  Future<List<Voice>> getAvailableVoices();

  /// Streams TTS preview updates
  Stream<Map<String, dynamic>> previewTts({
    required String text,
    required String voiceId,
    String? language,
    int emotion = 50,
    double speed = 1.0,
    double pitch = 1.0,
  });
}
