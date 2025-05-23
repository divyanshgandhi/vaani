import 'package:get_it/get_it.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:dio/dio.dart';
import 'package:vaani/core/network/api_client.dart';
import 'package:vaani/data/repositories/auth_repository_impl.dart';
import 'package:vaani/data/repositories/project_repository_impl.dart';
import 'package:vaani/data/repositories/tts_repository_impl.dart';
import 'package:vaani/domain/repositories/auth_repository.dart';
import 'package:vaani/domain/repositories/project_repository.dart';
import 'package:vaani/domain/repositories/tts_repository.dart';
import 'package:vaani/presentation/blocs/auth/auth_bloc.dart';
import 'package:vaani/presentation/blocs/project/project_bloc.dart';
import 'package:vaani/presentation/blocs/theme/theme_bloc.dart';
import 'package:vaani/presentation/blocs/language/language_bloc.dart';
import 'package:vaani/presentation/blocs/tts/tts_bloc.dart';

final GetIt getIt = GetIt.instance;

Future<void> initServiceLocator() async {
  // External dependencies
  final sharedPreferences = await SharedPreferences.getInstance();
  getIt.registerSingleton<SharedPreferences>(sharedPreferences);

  // Network
  getIt.registerLazySingleton<Dio>(() => Dio());
  getIt.registerLazySingleton<ApiClient>(() => ApiClient(getIt<Dio>()));

  // Repositories
  getIt.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(getIt<ApiClient>(), getIt<SharedPreferences>()),
  );

  getIt.registerLazySingleton<ProjectRepository>(
    () => ProjectRepositoryImpl(getIt<ApiClient>()),
  );

  getIt.registerLazySingleton<TtsRepository>(
    () => TtsRepositoryImpl(getIt<ApiClient>()),
  );

  // Initialize auth interceptor
  final apiClient = getIt<ApiClient>();
  final authRepository = getIt<AuthRepository>();
  apiClient.initAuthInterceptor(authRepository);

  // Blocs
  getIt.registerFactory<AuthBloc>(
    () => AuthBloc(authRepository: getIt<AuthRepository>()),
  );

  getIt.registerFactory<ThemeBloc>(
    () => ThemeBloc(),
  );

  getIt.registerFactory<LanguageBloc>(
    () => LanguageBloc(),
  );

  getIt.registerFactory<ProjectBloc>(
    () => ProjectBloc(projectRepository: getIt<ProjectRepository>()),
  );

  getIt.registerFactory<TtsBloc>(
    () => TtsBloc(ttsRepository: getIt<TtsRepository>()),
  );
}
