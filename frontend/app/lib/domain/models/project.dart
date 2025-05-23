class Project {
  final String id;
  final String title;
  final String description;
  final String audioUrl;
  final DateTime createdAt;
  final DateTime updatedAt;
  final bool isCompleted;

  Project({
    required this.id,
    required this.title,
    required this.description,
    this.audioUrl = '',
    required this.createdAt,
    required this.updatedAt,
    this.isCompleted = false,
  });

  factory Project.fromJson(Map<String, dynamic> json) {
    return Project(
      id: json['id'] ?? '',
      title: json['title'] ?? 'Untitled Project',
      description: json['description'] ?? '',
      audioUrl: json['audio_url'] ?? '',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'])
          : DateTime.now(),
      isCompleted: json['is_completed'] ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'audio_url': audioUrl,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      'is_completed': isCompleted,
    };
  }

  Project copyWith({
    String? id,
    String? title,
    String? description,
    String? audioUrl,
    DateTime? createdAt,
    DateTime? updatedAt,
    bool? isCompleted,
  }) {
    return Project(
      id: id ?? this.id,
      title: title ?? this.title,
      description: description ?? this.description,
      audioUrl: audioUrl ?? this.audioUrl,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      isCompleted: isCompleted ?? this.isCompleted,
    );
  }
}
