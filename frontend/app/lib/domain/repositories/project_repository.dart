abstract class ProjectRepository {
  Future<List<dynamic>> getProjects();
  Future<dynamic> getProjectById(String id);
  Future<dynamic> createProject(Map<String, dynamic> data);
  Future<dynamic> updateProject(String id, Map<String, dynamic> data);
  Future<bool> deleteProject(String id);
  Future<dynamic> generateAudio(
      String projectId, String text, Map<String, dynamic> options);
  Future<String> exportAudio(
      String projectId, Map<String, dynamic> exportOptions);
}
