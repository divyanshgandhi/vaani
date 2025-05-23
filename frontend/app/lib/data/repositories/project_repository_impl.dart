import 'package:vaani/core/network/api_client.dart';
import 'package:vaani/domain/repositories/project_repository.dart';

class ProjectRepositoryImpl implements ProjectRepository {
  final ApiClient _apiClient;

  ProjectRepositoryImpl(this._apiClient);

  @override
  Future<List<dynamic>> getProjects() async {
    try {
      final response = await _apiClient.get('/projects');
      if (response.statusCode == 200) {
        return response.data['projects'] ?? [];
      }
      return [];
    } catch (e) {
      return [];
    }
  }

  @override
  Future<dynamic> getProjectById(String id) async {
    try {
      final response = await _apiClient.get('/projects/$id');
      if (response.statusCode == 200) {
        return response.data;
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  @override
  Future<dynamic> createProject(Map<String, dynamic> data) async {
    try {
      final response = await _apiClient.post('/projects', data: data);
      if (response.statusCode == 201) {
        return response.data;
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  @override
  Future<dynamic> updateProject(String id, Map<String, dynamic> data) async {
    try {
      final response = await _apiClient.put('/projects/$id', data: data);
      if (response.statusCode == 200) {
        return response.data;
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  @override
  Future<bool> deleteProject(String id) async {
    try {
      final response = await _apiClient.delete('/projects/$id');
      return response.statusCode == 204;
    } catch (e) {
      return false;
    }
  }

  @override
  Future<dynamic> generateAudio(
      String projectId, String text, Map<String, dynamic> options) async {
    try {
      final response =
          await _apiClient.post('/projects/$projectId/generate', data: {
        'text': text,
        ...options,
      });
      if (response.statusCode == 200) {
        return response.data;
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  @override
  Future<String> exportAudio(
      String projectId, Map<String, dynamic> exportOptions) async {
    try {
      final response = await _apiClient.post('/projects/$projectId/export',
          data: exportOptions);
      if (response.statusCode == 200 && response.data['url'] != null) {
        return response.data['url'];
      }
      return '';
    } catch (e) {
      return '';
    }
  }
}
