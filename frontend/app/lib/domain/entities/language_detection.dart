import 'package:equatable/equatable.dart';

class LanguageDetection extends Equatable {
  final String language;
  final bool isIndic;

  const LanguageDetection({
    required this.language,
    required this.isIndic,
  });

  factory LanguageDetection.fromJson(Map<String, dynamic> json) {
    return LanguageDetection(
      language: json['language'] as String,
      isIndic: json['is_indic'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'language': language,
      'is_indic': isIndic,
    };
  }

  @override
  List<Object?> get props => [language, isIndic];
}
