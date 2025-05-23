import 'package:equatable/equatable.dart';

class TtsJob extends Equatable {
  final String id;
  final String status;
  final int? progress;
  final String? outputUrl;
  final String? errorMessage;
  final DateTime createdAt;
  final DateTime? completedAt;
  final double? duration;
  final String? message;

  const TtsJob({
    required this.id,
    required this.status,
    this.progress,
    this.outputUrl,
    this.errorMessage,
    required this.createdAt,
    this.completedAt,
    this.duration,
    this.message,
  });

  factory TtsJob.fromJson(Map<String, dynamic> json) {
    return TtsJob(
      id: json['job_id'] as String,
      status: json['status'] as String,
      progress: json['progress'] as int?,
      outputUrl: json['output_url'] as String?,
      errorMessage: json['error_message'] as String?,
      createdAt: DateTime.parse(json['created_at'] ??
          json['queue_time'] ??
          DateTime.now().toIso8601String()),
      completedAt: json['completed_at'] != null
          ? DateTime.parse(json['completed_at'])
          : null,
      duration: json['duration']?.toDouble(),
      message: json['message'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'job_id': id,
      'status': status,
      'progress': progress,
      'output_url': outputUrl,
      'error_message': errorMessage,
      'created_at': createdAt.toIso8601String(),
      'completed_at': completedAt?.toIso8601String(),
      'duration': duration,
      'message': message,
    };
  }

  TtsJob copyWith({
    String? id,
    String? status,
    int? progress,
    String? outputUrl,
    String? errorMessage,
    DateTime? createdAt,
    DateTime? completedAt,
    double? duration,
    String? message,
  }) {
    return TtsJob(
      id: id ?? this.id,
      status: status ?? this.status,
      progress: progress ?? this.progress,
      outputUrl: outputUrl ?? this.outputUrl,
      errorMessage: errorMessage ?? this.errorMessage,
      createdAt: createdAt ?? this.createdAt,
      completedAt: completedAt ?? this.completedAt,
      duration: duration ?? this.duration,
      message: message ?? this.message,
    );
  }

  @override
  List<Object?> get props => [
        id,
        status,
        progress,
        outputUrl,
        errorMessage,
        createdAt,
        completedAt,
        duration,
        message,
      ];
}
