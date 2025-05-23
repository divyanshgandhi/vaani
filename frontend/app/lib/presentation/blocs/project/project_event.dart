import 'package:equatable/equatable.dart';

abstract class ProjectEvent extends Equatable {
  const ProjectEvent();

  @override
  List<Object> get props => [];
}

class ProjectsLoaded extends ProjectEvent {
  const ProjectsLoaded();
}

class ProjectRefreshed extends ProjectEvent {
  const ProjectRefreshed();
}

class ProjectCreated extends ProjectEvent {
  final String title;
  final String description;

  const ProjectCreated({
    required this.title,
    required this.description,
  });

  @override
  List<Object> get props => [title, description];
}

class ProjectDeleted extends ProjectEvent {
  final String id;

  const ProjectDeleted({required this.id});

  @override
  List<Object> get props => [id];
}
