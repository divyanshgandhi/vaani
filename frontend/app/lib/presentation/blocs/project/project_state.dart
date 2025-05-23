import 'package:equatable/equatable.dart';
import 'package:vaani/domain/models/project.dart';

abstract class ProjectState extends Equatable {
  const ProjectState();

  @override
  List<Object> get props => [];
}

class ProjectInitial extends ProjectState {
  const ProjectInitial();
}

class ProjectLoading extends ProjectState {
  const ProjectLoading();
}

class ProjectLoaded extends ProjectState {
  final List<Project> projects;

  const ProjectLoaded({required this.projects});

  @override
  List<Object> get props => [projects];

  ProjectLoaded copyWith({
    List<Project>? projects,
  }) {
    return ProjectLoaded(
      projects: projects ?? this.projects,
    );
  }
}

class ProjectError extends ProjectState {
  final String message;

  const ProjectError({required this.message});

  @override
  List<Object> get props => [message];
}

class ProjectCreationSuccess extends ProjectState {
  final Project project;

  const ProjectCreationSuccess({required this.project});

  @override
  List<Object> get props => [project];
}

class ProjectCreationFailure extends ProjectState {
  final String message;

  const ProjectCreationFailure({required this.message});

  @override
  List<Object> get props => [message];
}
