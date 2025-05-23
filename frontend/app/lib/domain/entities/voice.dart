import 'package:equatable/equatable.dart';

class Voice extends Equatable {
  final String id;
  final String name;
  final String language;
  final String gender;
  final String? previewUrl;

  const Voice({
    required this.id,
    required this.name,
    required this.language,
    required this.gender,
    this.previewUrl,
  });

  factory Voice.fromJson(Map<String, dynamic> json) {
    return Voice(
      id: json['id'] as String,
      name: json['name'] as String,
      language: json['language'] as String,
      gender: json['gender'] as String,
      previewUrl: json['preview_url'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'language': language,
      'gender': gender,
      'preview_url': previewUrl,
    };
  }

  @override
  List<Object?> get props => [id, name, language, gender, previewUrl];
}
