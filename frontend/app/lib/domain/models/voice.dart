import 'package:equatable/equatable.dart';

enum VoiceType {
  stock,
  cloned,
}

class Voice extends Equatable {
  final String id;
  final String name;
  final String language;
  final String avatarUrl;
  final String sampleAudioUrl;
  final VoiceType type;

  const Voice({
    required this.id,
    required this.name,
    required this.language,
    required this.avatarUrl,
    required this.sampleAudioUrl,
    required this.type,
  });

  @override
  List<Object?> get props => [
        id,
        name,
        language,
        avatarUrl,
        sampleAudioUrl,
        type,
      ];

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'language': language,
      'avatarUrl': avatarUrl,
      'sampleAudioUrl': sampleAudioUrl,
      'type': type.toString(),
    };
  }

  factory Voice.fromJson(Map<String, dynamic> json) {
    return Voice(
      id: json['id'] as String,
      name: json['name'] as String,
      language: json['language'] as String,
      avatarUrl: json['avatarUrl'] as String,
      sampleAudioUrl: json['sampleAudioUrl'] as String,
      type: json['type'] == 'VoiceType.cloned'
          ? VoiceType.cloned
          : VoiceType.stock,
    );
  }
}
