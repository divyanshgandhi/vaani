import 'package:equatable/equatable.dart';

abstract class TtsEvent extends Equatable {
  const TtsEvent();

  @override
  List<Object?> get props => [];
}

class TtsLanguageDetectionRequested extends TtsEvent {
  final String text;

  const TtsLanguageDetectionRequested(this.text);

  @override
  List<Object?> get props => [text];
}

class TtsGenerationRequested extends TtsEvent {
  final String text;
  final String voiceId;
  final String? language;
  final int emotion;
  final double speed;
  final double pitch;
  final String outputType;

  const TtsGenerationRequested({
    required this.text,
    required this.voiceId,
    this.language,
    this.emotion = 50,
    this.speed = 1.0,
    this.pitch = 1.0,
    this.outputType = 'mp3',
  });

  @override
  List<Object?> get props => [
        text,
        voiceId,
        language,
        emotion,
        speed,
        pitch,
        outputType,
      ];
}

class TtsJobStatusRequested extends TtsEvent {
  final String jobId;

  const TtsJobStatusRequested(this.jobId);

  @override
  List<Object?> get props => [jobId];
}

class TtsVoicesRequested extends TtsEvent {
  const TtsVoicesRequested();
}

class TtsPreviewRequested extends TtsEvent {
  final String text;
  final String voiceId;
  final String? language;
  final int emotion;
  final double speed;
  final double pitch;

  const TtsPreviewRequested({
    required this.text,
    required this.voiceId,
    this.language,
    this.emotion = 50,
    this.speed = 1.0,
    this.pitch = 1.0,
  });

  @override
  List<Object?> get props => [
        text,
        voiceId,
        language,
        emotion,
        speed,
        pitch,
      ];
}

class TtsPreviewStopped extends TtsEvent {
  const TtsPreviewStopped();
}

class TtsReset extends TtsEvent {
  const TtsReset();
}
