import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:vaani/domain/repositories/tts_repository.dart';
import 'package:vaani/presentation/blocs/tts/tts_event.dart';
import 'package:vaani/presentation/blocs/tts/tts_state.dart';

class TtsBloc extends Bloc<TtsEvent, TtsState> {
  final TtsRepository _ttsRepository;
  StreamSubscription? _previewSubscription;

  TtsBloc({required TtsRepository ttsRepository})
      : _ttsRepository = ttsRepository,
        super(const TtsInitial()) {
    on<TtsLanguageDetectionRequested>(_onLanguageDetectionRequested);
    on<TtsGenerationRequested>(_onGenerationRequested);
    on<TtsJobStatusRequested>(_onJobStatusRequested);
    on<TtsVoicesRequested>(_onVoicesRequested);
    on<TtsPreviewRequested>(_onPreviewRequested);
    on<TtsPreviewStopped>(_onPreviewStopped);
    on<TtsReset>(_onReset);
  }

  Future<void> _onLanguageDetectionRequested(
    TtsLanguageDetectionRequested event,
    Emitter<TtsState> emit,
  ) async {
    emit(const TtsLanguageDetectionLoading());
    try {
      final languageDetection = await _ttsRepository.detectLanguage(event.text);
      emit(TtsLanguageDetectionSuccess(languageDetection));
    } catch (e) {
      emit(TtsLanguageDetectionFailure(e.toString()));
    }
  }

  Future<void> _onGenerationRequested(
    TtsGenerationRequested event,
    Emitter<TtsState> emit,
  ) async {
    emit(const TtsGenerationLoading());
    try {
      final job = await _ttsRepository.generateTts(
        text: event.text,
        voiceId: event.voiceId,
        language: event.language,
        emotion: event.emotion,
        speed: event.speed,
        pitch: event.pitch,
        outputType: event.outputType,
      );
      emit(TtsGenerationSuccess(job));
    } catch (e) {
      emit(TtsGenerationFailure(e.toString()));
    }
  }

  Future<void> _onJobStatusRequested(
    TtsJobStatusRequested event,
    Emitter<TtsState> emit,
  ) async {
    emit(TtsJobStatusLoading(event.jobId));
    try {
      final job = await _ttsRepository.getJobStatus(event.jobId);
      emit(TtsJobStatusSuccess(job));
    } catch (e) {
      emit(TtsJobStatusFailure(e.toString()));
    }
  }

  Future<void> _onVoicesRequested(
    TtsVoicesRequested event,
    Emitter<TtsState> emit,
  ) async {
    emit(const TtsVoicesLoading());
    try {
      final voices = await _ttsRepository.getAvailableVoices();
      emit(TtsVoicesSuccess(voices));
    } catch (e) {
      emit(TtsVoicesFailure(e.toString()));
    }
  }

  Future<void> _onPreviewRequested(
    TtsPreviewRequested event,
    Emitter<TtsState> emit,
  ) async {
    // Cancel any existing preview
    await _previewSubscription?.cancel();

    emit(const TtsPreviewStarted());

    try {
      _previewSubscription = _ttsRepository
          .previewTts(
        text: event.text,
        voiceId: event.voiceId,
        language: event.language,
        emotion: event.emotion,
        speed: event.speed,
        pitch: event.pitch,
      )
          .listen(
        (data) {
          switch (data['type']) {
            case 'job_created':
            case 'status_update':
              emit(TtsPreviewProgress(data));
              break;
            case 'completed':
              emit(TtsPreviewCompleted(
                outputUrl: data['output_url'],
                jobId: data['job_id'],
                duration: data['duration']?.toDouble(),
              ));
              break;
            case 'error':
              emit(TtsPreviewFailure(data['error']));
              break;
          }
        },
        onError: (error) {
          emit(TtsPreviewFailure(error.toString()));
        },
      );
    } catch (e) {
      emit(TtsPreviewFailure(e.toString()));
    }
  }

  Future<void> _onPreviewStopped(
    TtsPreviewStopped event,
    Emitter<TtsState> emit,
  ) async {
    await _previewSubscription?.cancel();
    _previewSubscription = null;
    emit(const TtsInitial());
  }

  Future<void> _onReset(
    TtsReset event,
    Emitter<TtsState> emit,
  ) async {
    await _previewSubscription?.cancel();
    _previewSubscription = null;
    emit(const TtsInitial());
  }

  @override
  Future<void> close() {
    _previewSubscription?.cancel();
    return super.close();
  }
}
