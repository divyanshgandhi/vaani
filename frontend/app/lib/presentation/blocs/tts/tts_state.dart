import 'package:equatable/equatable.dart';
import 'package:vaani/domain/entities/language_detection.dart';
import 'package:vaani/domain/entities/tts_job.dart';
import 'package:vaani/domain/entities/voice.dart';

abstract class TtsState extends Equatable {
  const TtsState();

  @override
  List<Object?> get props => [];
}

class TtsInitial extends TtsState {
  const TtsInitial();
}

class TtsLoading extends TtsState {
  const TtsLoading();
}

// Language Detection States
class TtsLanguageDetectionLoading extends TtsState {
  const TtsLanguageDetectionLoading();
}

class TtsLanguageDetectionSuccess extends TtsState {
  final LanguageDetection languageDetection;

  const TtsLanguageDetectionSuccess(this.languageDetection);

  @override
  List<Object?> get props => [languageDetection];
}

class TtsLanguageDetectionFailure extends TtsState {
  final String error;

  const TtsLanguageDetectionFailure(this.error);

  @override
  List<Object?> get props => [error];
}

// TTS Generation States
class TtsGenerationLoading extends TtsState {
  const TtsGenerationLoading();
}

class TtsGenerationSuccess extends TtsState {
  final TtsJob job;

  const TtsGenerationSuccess(this.job);

  @override
  List<Object?> get props => [job];
}

class TtsGenerationFailure extends TtsState {
  final String error;

  const TtsGenerationFailure(this.error);

  @override
  List<Object?> get props => [error];
}

// Job Status States
class TtsJobStatusLoading extends TtsState {
  final String jobId;

  const TtsJobStatusLoading(this.jobId);

  @override
  List<Object?> get props => [jobId];
}

class TtsJobStatusSuccess extends TtsState {
  final TtsJob job;

  const TtsJobStatusSuccess(this.job);

  @override
  List<Object?> get props => [job];
}

class TtsJobStatusFailure extends TtsState {
  final String error;

  const TtsJobStatusFailure(this.error);

  @override
  List<Object?> get props => [error];
}

// Voices States
class TtsVoicesLoading extends TtsState {
  const TtsVoicesLoading();
}

class TtsVoicesSuccess extends TtsState {
  final List<Voice> voices;

  const TtsVoicesSuccess(this.voices);

  @override
  List<Object?> get props => [voices];
}

class TtsVoicesFailure extends TtsState {
  final String error;

  const TtsVoicesFailure(this.error);

  @override
  List<Object?> get props => [error];
}

// Preview States
class TtsPreviewStarted extends TtsState {
  const TtsPreviewStarted();
}

class TtsPreviewProgress extends TtsState {
  final Map<String, dynamic> data;

  const TtsPreviewProgress(this.data);

  @override
  List<Object?> get props => [data];
}

class TtsPreviewCompleted extends TtsState {
  final String outputUrl;
  final String jobId;
  final double? duration;

  const TtsPreviewCompleted({
    required this.outputUrl,
    required this.jobId,
    this.duration,
  });

  @override
  List<Object?> get props => [outputUrl, jobId, duration];
}

class TtsPreviewFailure extends TtsState {
  final String error;

  const TtsPreviewFailure(this.error);

  @override
  List<Object?> get props => [error];
}
