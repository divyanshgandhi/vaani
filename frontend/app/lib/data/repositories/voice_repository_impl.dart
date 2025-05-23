import 'package:dio/dio.dart';
import '../../domain/models/voice.dart';
import '../../domain/repositories/voice_repository.dart';

class VoiceRepositoryImpl implements VoiceRepository {
  final Dio _dio;

  VoiceRepositoryImpl(this._dio);

  @override
  Future<List<Voice>> getVoices() async {
    try {
      final response = await _dio.get('/voices');
      final List<dynamic> voicesJson = response.data['voices'] ?? [];
      return voicesJson.map((json) => Voice.fromJson(json)).toList();
    } catch (e) {
      // Return some mock voices if API fails
      return [
        const Voice(
          id: 'en-female-1',
          name: 'Sophia',
          language: 'en',
          avatarUrl: 'assets/images/avatar_f1.png',
          sampleAudioUrl: 'assets/audio/sample_f1.mp3',
          type: VoiceType.stock,
        ),
        const Voice(
          id: 'en-male-1',
          name: 'James',
          language: 'en',
          avatarUrl: 'assets/images/avatar_m1.png',
          sampleAudioUrl: 'assets/audio/sample_m1.mp3',
          type: VoiceType.stock,
        ),
        const Voice(
          id: 'hi-female-1',
          name: 'Priya',
          language: 'hi',
          avatarUrl: 'assets/images/avatar_f2.png',
          sampleAudioUrl: 'assets/audio/sample_f2.mp3',
          type: VoiceType.stock,
        ),
      ];
    }
  }

  @override
  Future<Voice?> getVoice(String id) async {
    try {
      final response = await _dio.get('/voices/$id');
      return Voice.fromJson(response.data);
    } catch (e) {
      // Lookup in local list
      final voices = await getVoices();
      return voices.firstWhere(
        (voice) => voice.id == id,
        orElse: () => voices.first,
      );
    }
  }
}
