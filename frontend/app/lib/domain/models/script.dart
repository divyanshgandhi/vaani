import 'package:equatable/equatable.dart';
import 'package:flutter/foundation.dart';

class Script extends Equatable {
  final String text;
  final String voiceId;
  final String voiceName;
  final String language;
  final int emotionLevel;
  final int charCount;

  const Script({
    this.text = '',
    this.voiceId = 'default',
    this.voiceName = 'Default Voice',
    this.language = 'en',
    this.emotionLevel = 50,
    this.charCount = 0,
  });

  Script copyWith({
    String? text,
    String? voiceId,
    String? voiceName,
    String? language,
    int? emotionLevel,
    int? charCount,
  }) {
    return Script(
      text: text ?? this.text,
      voiceId: voiceId ?? this.voiceId,
      voiceName: voiceName ?? this.voiceName,
      language: language ?? this.language,
      emotionLevel: emotionLevel ?? this.emotionLevel,
      charCount: charCount ?? this.charCount,
    );
  }

  @override
  List<Object?> get props => [
        text,
        voiceId,
        voiceName,
        language,
        emotionLevel,
        charCount,
      ];

  Map<String, dynamic> toJson() {
    return {
      'text': text,
      'voiceId': voiceId,
      'voiceName': voiceName,
      'language': language,
      'emotionLevel': emotionLevel,
      'charCount': charCount,
    };
  }

  factory Script.fromJson(Map<String, dynamic> json) {
    return Script(
      text: json['text'] ?? '',
      voiceId: json['voiceId'] ?? 'default',
      voiceName: json['voiceName'] ?? 'Default Voice',
      language: json['language'] ?? 'en',
      emotionLevel: json['emotionLevel'] ?? 50,
      charCount: json['charCount'] ?? 0,
    );
  }
}
