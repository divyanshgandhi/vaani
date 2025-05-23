import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';
import 'package:vaani/core/theme/app_theme.dart';
import 'package:vaani/presentation/blocs/auth/auth_bloc.dart';
import 'package:vaani/presentation/blocs/project/project_bloc.dart';
import 'package:vaani/presentation/blocs/project/project_event.dart';
import 'package:vaani/presentation/blocs/theme/theme_bloc.dart';
import 'package:vaani/presentation/blocs/language/language_bloc.dart';
import 'package:vaani/presentation/features/auth/phone_input_screen.dart';
import 'package:vaani/presentation/features/editor/editor_screen.dart';
import 'package:vaani/presentation/features/home/home_screen.dart';
import 'package:vaani/presentation/features/splash/splash_screen.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class VaaniApp extends StatefulWidget {
  const VaaniApp({super.key});

  @override
  State<VaaniApp> createState() => _VaaniAppState();
}

class _VaaniAppState extends State<VaaniApp> {
  late final GoRouter _router;

  @override
  void initState() {
    super.initState();

    _router = GoRouter(
      initialLocation: '/splash',
      routes: [
        GoRoute(
          path: '/splash',
          builder: (context, state) => const SplashScreen(),
        ),
        GoRoute(
          path: '/login',
          builder: (context, state) => const PhoneInputScreen(),
        ),
        GoRoute(
          path: '/home',
          builder: (context, state) => const HomeScreen(),
        ),
        GoRoute(
          path: '/editor',
          builder: (context, state) => const EditorScreen(),
        ),
        GoRoute(
          path: '/editor/:projectId',
          builder: (context, state) {
            final projectId = state.pathParameters['projectId'];
            return EditorScreen(projectId: projectId);
          },
        ),
      ],
      redirect: (context, state) {
        final authState = context.read<AuthBloc>().state;
        final isLoginRoute = state.matchedLocation == '/login';
        final isSplashRoute = state.matchedLocation == '/splash';

        // If we're at the splash screen, don't redirect
        if (isSplashRoute) {
          return null;
        }

        // If user is authenticated but on login page, redirect to home
        if (authState is AuthSuccess && isLoginRoute) {
          return '/home';
        }

        // If user is not authenticated and not on login page, redirect to login
        if (authState is! AuthSuccess && !isLoginRoute) {
          return '/login';
        }

        return null;
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(
          create: (context) =>
              GetIt.instance<AuthBloc>()..add(const AuthStarted()),
        ),
        BlocProvider(
          create: (context) => GetIt.instance<ThemeBloc>(),
        ),
        BlocProvider(
          create: (context) => GetIt.instance<LanguageBloc>(),
        ),
        BlocProvider(
          create: (context) =>
              GetIt.instance<ProjectBloc>()..add(const ProjectsLoaded()),
        ),
      ],
      child: Builder(
        builder: (context) {
          final themeState = context.watch<ThemeBloc>().state;
          final languageState = context.watch<LanguageBloc>().state;

          return MaterialApp.router(
            title: 'Vaani',
            theme: AppTheme.lightTheme,
            darkTheme: AppTheme.darkTheme,
            themeMode: themeState.themeMode,
            locale: languageState.locale,
            supportedLocales: const [
              Locale('en', ''), // English
              Locale('hi', ''), // Hindi
            ],
            localizationsDelegates: const [
              AppLocalizations.delegate,
              GlobalMaterialLocalizations.delegate,
              GlobalWidgetsLocalizations.delegate,
              GlobalCupertinoLocalizations.delegate,
            ],
            routerConfig: _router,
          );
        },
      ),
    );
  }
}
