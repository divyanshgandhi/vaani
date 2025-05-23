import 'package:dio/dio.dart';
import 'package:vaani/core/network/api_client.dart';
import 'package:vaani/domain/entities/tts_job.dart';
import 'package:vaani/domain/entities/language_detection.dart';
import 'package:vaani/domain/entities/voice.dart';
import 'package:vaani/domain/repositories/tts_repository.dart';

class TtsRepositoryImpl implements TtsRepository {
  final ApiClient _apiClient;

  TtsRepositoryImpl(this._apiClient);

  @override
  Future<LanguageDetection> detectLanguage(String text) async {
    try {
      final response = await _apiClient.post('/detect-language', data: {
        'text': text,
      });

      return LanguageDetection.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  @override
  Future<TtsJob> generateTts({
    required String text,
    required String voiceId,
    String? language,
    int emotion = 50,
    double speed = 1.0,
    double pitch = 1.0,
    String outputType = 'mp3',
  }) async {
    try {
      final response = await _apiClient.post('/generate', data: {
        'text': text,
        'voice_id': voiceId,
        'language': language,
        'emotion': emotion,
        'speed': speed,
        'pitch': pitch,
        'output_type': outputType,
      });

      return TtsJob.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  @override
  Future<TtsJob> getJobStatus(String jobId) async {
    try {
      final response = await _apiClient.get('/jobs/$jobId');
      return TtsJob.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  @override
  Future<String> getDownloadUrl(String jobId) async {
    try {
      final response = await _apiClient.get('/download/$jobId');
      // The backend redirects to the actual download URL
      // We can extract it from the response or location header
      return response.realUri.toString();
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  @override
  Future<List<Voice>> getAvailableVoices() async {
    try {
      final response = await _apiClient.get('/voices');
      final List<dynamic> data = response.data;
      return data.map((json) => Voice.fromJson(json)).toList();
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  @override
  Stream<Map<String, dynamic>> previewTts({
    required String text,
    required String voiceId,
    String? language,
    int emotion = 50,
    double speed = 1.0,
    double pitch = 1.0,
  }) async* {
    // This would use WebSocket connection to the /preview endpoint
    // For now, we'll implement a basic stream that polls the job status

    try {
      // First, generate the TTS job
      final job = await generateTts(
        text: text,
        voiceId: voiceId,
        language: language,
        emotion: emotion,
        speed: speed,
        pitch: pitch,
      );

      // Yield initial job creation
      yield {
        'type': 'job_created',
        'job_id': job.id,
        'status': job.status,
      };

      // Poll for updates until completion
      while (job.status != 'complete' && job.status != 'failed') {
        await Future.delayed(const Duration(seconds: 2));

        final updatedJob = await getJobStatus(job.id);

        yield {
          'type': 'status_update',
          'job_id': updatedJob.id,
          'status': updatedJob.status,
          'progress': updatedJob.progress,
        };

        if (updatedJob.status == 'complete') {
          yield {
            'type': 'completed',
            'job_id': updatedJob.id,
            'output_url': updatedJob.outputUrl,
            'duration': updatedJob.duration,
          };
          break;
        } else if (updatedJob.status == 'failed') {
          yield {
            'type': 'error',
            'job_id': updatedJob.id,
            'error': updatedJob.errorMessage ?? 'Unknown error',
          };
          break;
        }
      }
    } catch (e) {
      yield {
        'type': 'error',
        'error': e.toString(),
      };
    }
  }

  Exception _handleError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
        return Exception(
            'Connection timeout. Please check your internet connection.');
      case DioExceptionType.badResponse:
        final statusCode = e.response?.statusCode;
        final message = e.response?.data?['message'] ?? 'Server error';
        return Exception('Server error ($statusCode): $message');
      default:
        return Exception('Network error: ${e.message}');
    }
  }
}
