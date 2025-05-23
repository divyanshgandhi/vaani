import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:vaani/domain/models/project.dart';
import 'package:vaani/domain/repositories/project_repository.dart';
import 'package:vaani/presentation/blocs/project/project_event.dart';
import 'package:vaani/presentation/blocs/project/project_state.dart';

class ProjectBloc extends Bloc<ProjectEvent, ProjectState> {
  final ProjectRepository _projectRepository;

  ProjectBloc({required ProjectRepository projectRepository})
      : _projectRepository = projectRepository,
        super(const ProjectInitial()) {
    on<ProjectsLoaded>(_onProjectsLoaded);
    on<ProjectRefreshed>(_onProjectRefreshed);
    on<ProjectCreated>(_onProjectCreated);
    on<ProjectDeleted>(_onProjectDeleted);
  }

  Future<void> _onProjectsLoaded(
    ProjectsLoaded event,
    Emitter<ProjectState> emit,
  ) async {
    emit(const ProjectLoading());
    try {
      final projectsData = await _projectRepository.getProjects();
      final projects = projectsData
          .map((projectJson) => Project.fromJson(projectJson))
          .toList();
      emit(ProjectLoaded(projects: projects));
    } catch (e) {
      emit(ProjectError(message: e.toString()));
    }
  }

  Future<void> _onProjectRefreshed(
    ProjectRefreshed event,
    Emitter<ProjectState> emit,
  ) async {
    try {
      final projectsData = await _projectRepository.getProjects();
      final projects = projectsData
          .map((projectJson) => Project.fromJson(projectJson))
          .toList();
      emit(ProjectLoaded(projects: projects));
    } catch (e) {
      // Keep the current state if refresh fails
      if (state is ProjectLoaded) {
        emit(state);
      } else {
        emit(ProjectError(message: e.toString()));
      }
    }
  }

  Future<void> _onProjectCreated(
    ProjectCreated event,
    Emitter<ProjectState> emit,
  ) async {
    try {
      final projectData = await _projectRepository.createProject({
        'title': event.title,
        'description': event.description,
      });

      if (projectData != null) {
        final project = Project.fromJson(projectData);
        emit(ProjectCreationSuccess(project: project));

        // Load projects again to refresh the list
        add(const ProjectsLoaded());
      } else {
        emit(const ProjectCreationFailure(message: 'Failed to create project'));
      }
    } catch (e) {
      emit(ProjectCreationFailure(message: e.toString()));
    }
  }

  Future<void> _onProjectDeleted(
    ProjectDeleted event,
    Emitter<ProjectState> emit,
  ) async {
    try {
      final success = await _projectRepository.deleteProject(event.id);

      if (success) {
        // If we have a loaded state, filter out the deleted project
        if (state is ProjectLoaded) {
          final currentProjects = (state as ProjectLoaded).projects;
          final updatedProjects = currentProjects
              .where((project) => project.id != event.id)
              .toList();
          emit(ProjectLoaded(projects: updatedProjects));
        }
      }
    } catch (e) {
      // Keep current state
      if (state is ProjectLoaded) {
        emit(state);
      } else {
        emit(ProjectError(message: e.toString()));
      }
    }
  }
}
