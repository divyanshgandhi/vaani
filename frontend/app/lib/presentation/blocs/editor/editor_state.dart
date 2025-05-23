import 'package:equatable/equatable.dart';
import '../../../domain/models/script.dart';
import '../../../domain/models/voice.dart';

enum EditorStatus {
  initial,
  loading,
  loadingVoices,
  ready,
  error,
}

class EditorState extends Equatable {
  final EditorStatus status;
  final Script script;
  final List<Voice> availableVoices;
  final double pricePerChar;
  final double totalCost;
  final String? errorMessage;

  const EditorState({
    this.status = EditorStatus.initial,
    this.script = const Script(),
    this.availableVoices = const [],
    this.pricePerChar = 0.01,
    this.totalCost = 0.0,
    this.errorMessage,
  });

  bool get isOverCharLimit => script.charCount > 1000;

  EditorState copyWith({
    EditorStatus? status,
    Script? script,
    List<Voice>? availableVoices,
    double? pricePerChar,
    double? totalCost,
    String? errorMessage,
  }) {
    return EditorState(
      status: status ?? this.status,
      script: script ?? this.script,
      availableVoices: availableVoices ?? this.availableVoices,
      pricePerChar: pricePerChar ?? this.pricePerChar,
      totalCost: totalCost ?? this.totalCost,
      errorMessage: errorMessage,
    );
  }

  @override
  List<Object?> get props => [
        status,
        script,
        availableVoices,
        pricePerChar,
        totalCost,
        errorMessage,
      ];
}
