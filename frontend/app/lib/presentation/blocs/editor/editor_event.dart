import 'package:equatable/equatable.dart';

abstract class EditorEvent extends Equatable {
  const EditorEvent();

  @override
  List<Object?> get props => [];
}

class TextChanged extends EditorEvent {
  final String text;

  const TextChanged(this.text);

  @override
  List<Object?> get props => [text];
}

class ApplyAutoPunctuation extends EditorEvent {}

class VoiceSelected extends EditorEvent {
  final String voiceId;
  final String voiceName;
  final String language;

  const VoiceSelected({
    required this.voiceId,
    required this.voiceName,
    required this.language,
  });

  @override
  List<Object?> get props => [voiceId, voiceName, language];
}

class EmotionChanged extends EditorEvent {
  final int level;

  const EmotionChanged(this.level);

  @override
  List<Object?> get props => [level];
}

class FetchVoices extends EditorEvent {}

class FetchPricing extends EditorEvent {}
