import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:go_router/go_router.dart';
import 'package:vaani/core/components/modern_button.dart';
import 'package:vaani/core/components/modern_card.dart';
import 'package:vaani/core/theme/app_theme.dart';
import 'package:vaani/core/animations/app_animations.dart';
import 'package:vaani/domain/models/project.dart';
import 'package:vaani/presentation/blocs/auth/auth_bloc.dart';
import 'package:vaani/presentation/blocs/project/project_bloc.dart';
import 'package:vaani/presentation/blocs/project/project_event.dart';
import 'package:vaani/presentation/blocs/project/project_state.dart';
import 'package:vaani/presentation/blocs/theme/theme_bloc.dart';
import 'package:vaani/presentation/features/home/create_project_dialog.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> with TickerProviderStateMixin {
  late AnimationController _headerAnimationController;
  late AnimationController _fabAnimationController;

  @override
  void initState() {
    super.initState();
    _headerAnimationController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    _fabAnimationController = AnimationController(
      duration: const Duration(milliseconds: 600),
      vsync: this,
    );

    context.read<ProjectBloc>().add(const ProjectsLoaded());

    // Start animations
    _headerAnimationController.forward();
    Future.delayed(const Duration(milliseconds: 400), () {
      if (mounted) _fabAnimationController.forward();
    });
  }

  @override
  void dispose() {
    _headerAnimationController.dispose();
    _fabAnimationController.dispose();
    super.dispose();
  }

  void _showCreateProjectDialog(BuildContext context) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const CreateProjectDialog(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      backgroundColor: colorScheme.background,
      body: CustomScrollView(
        slivers: [
          // Modern App Bar
          SliverAppBar(
            expandedHeight: 120,
            floating: true,
            pinned: true,
            backgroundColor: Colors.transparent,
            elevation: 0,
            flexibleSpace: FlexibleSpaceBar(
              background: Container(
                decoration: BoxDecoration(
                  gradient: theme.brightness == Brightness.light
                      ? AppTheme.primaryGradient
                      : AppTheme.primaryGradientDark,
                ),
                child: SafeArea(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(20, 20, 20, 0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Welcome back',
                                  style: theme.textTheme.bodyLarge?.copyWith(
                                    color: Colors.white.withOpacity(0.9),
                                  ),
                                ),
                                Text(
                                  'My Projects',
                                  style:
                                      theme.textTheme.headlineMedium?.copyWith(
                                    color: Colors.white,
                                    fontWeight: FontWeight.w700,
                                  ),
                                ),
                              ],
                            )
                                .animate(controller: _headerAnimationController)
                                .fadeIn(duration: 600.ms)
                                .slideX(begin: -0.3, end: 0),
                            Row(
                              children: [
                                _buildThemeToggle(),
                                const SizedBox(width: 8),
                                _buildSettingsButton(),
                                const SizedBox(width: 8),
                                _buildLogoutButton(),
                              ],
                            )
                                .animate(controller: _headerAnimationController)
                                .fadeIn(duration: 600.ms, delay: 200.ms)
                                .slideX(begin: 0.3, end: 0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ),

          // Stats Section
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: _buildStatsSection()
                  .animate()
                  .fadeIn(
                    duration: 600.ms,
                    delay: 300.ms,
                  )
                  .slideY(begin: 0.3, end: 0),
            ),
          ),

          // Projects Grid
          BlocBuilder<ProjectBloc, ProjectState>(
            builder: (context, state) {
              if (state is ProjectLoading) {
                return SliverToBoxAdapter(
                  child: _buildLoadingState(),
                );
              } else if (state is ProjectError) {
                return SliverToBoxAdapter(
                  child: _buildErrorState(state.message, l10n),
                );
              } else if (state is ProjectLoaded) {
                final projects = state.projects;

                if (projects.isEmpty) {
                  return SliverToBoxAdapter(
                    child: _buildEmptyState(l10n),
                  );
                }

                return SliverPadding(
                  padding: const EdgeInsets.fromLTRB(20, 0, 20, 100),
                  sliver: SliverGrid(
                    delegate: SliverChildBuilderDelegate(
                      (context, index) {
                        final project = projects[index];
                        return ProjectCard(
                          title: project.title,
                          description: project.description,
                          lastModified: project.updatedAt,
                          hasAudio: project.audioUrl.isNotEmpty,
                          accentColor: _getProjectAccentColor(index),
                          onTap: () => context.push('/editor/${project.id}'),
                          onDelete: () =>
                              _showDeleteConfirmation(context, project),
                        )
                            .animate()
                            .fadeIn(
                              duration: 400.ms,
                              delay: Duration(milliseconds: 100 * index),
                            )
                            .slideY(
                                begin: 0.3, end: 0, curve: Curves.easeOutCubic);
                      },
                      childCount: projects.length,
                    ),
                    gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: _getCrossAxisCount(context),
                      childAspectRatio: 1.2,
                      crossAxisSpacing: 16,
                      mainAxisSpacing: 16,
                    ),
                  ),
                );
              }

              return SliverToBoxAdapter(
                child: Center(
                  child: Text(
                    'Something went wrong',
                    style: theme.textTheme.bodyLarge,
                  ),
                ),
              );
            },
          ),
        ],
      ),
      floatingActionButton: AnimatedBuilder(
        animation: _fabAnimationController,
        builder: (context, child) {
          return Transform.scale(
            scale: _fabAnimationController.value,
            child: ModernFloatingActionButton(
              onPressed: () => _showCreateProjectDialog(context),
              child: const Icon(
                Icons.add,
                color: Colors.white,
                size: 28,
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildThemeToggle() {
    return BlocBuilder<ThemeBloc, ThemeState>(
      builder: (context, state) {
        final isDark = state.themeMode == ThemeMode.dark;
        return GestureDetector(
          onTap: () {
            context.read<ThemeBloc>().add(ThemeToggled());
          },
          child: Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(
              isDark ? Icons.light_mode : Icons.dark_mode,
              color: Colors.white,
              size: 20,
            ),
          ),
        );
      },
    );
  }

  Widget _buildSettingsButton() {
    return GestureDetector(
      onTap: () {
        // TODO: Navigate to settings
      },
      child: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.2),
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Icon(
          Icons.settings,
          color: Colors.white,
          size: 20,
        ),
      ),
    );
  }

  Widget _buildLogoutButton() {
    return GestureDetector(
      onTap: () {
        context.read<AuthBloc>().add(const AuthLogoutRequested());
      },
      child: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.2),
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Icon(
          Icons.logout,
          color: Colors.white,
          size: 20,
        ),
      ),
    );
  }

  Widget _buildStatsSection() {
    return BlocBuilder<ProjectBloc, ProjectState>(
      builder: (context, state) {
        if (state is! ProjectLoaded) return const SizedBox.shrink();

        final projects = state.projects;
        final audioProjects =
            projects.where((p) => p.audioUrl.isNotEmpty).length;

        return Row(
          children: [
            Expanded(
              child: StatsCard(
                title: 'Total Projects',
                value: '${projects.length}',
                icon: Icons.folder_outlined,
                iconColor: Theme.of(context).colorScheme.primary,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: StatsCard(
                title: 'Audio Generated',
                value: '$audioProjects',
                icon: Icons.music_note_outlined,
                iconColor: Theme.of(context).colorScheme.secondary,
              ),
            ),
          ],
        );
      },
    );
  }

  Widget _buildLoadingState() {
    return Padding(
      padding: const EdgeInsets.all(40),
      child: Column(
        children: [
          AppAnimations.loadingSpinner(
            color: Theme.of(context).colorScheme.primary,
            size: 32,
          ),
          const SizedBox(height: 16),
          Text(
            'Loading your projects...',
            style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(String message, AppLocalizations l10n) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.all(40),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: theme.colorScheme.errorContainer,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Icon(
              Icons.error_outline,
              size: 48,
              color: theme.colorScheme.onErrorContainer,
            ),
          ),
          const SizedBox(height: 16),
          Text(
            message,
            style: theme.textTheme.bodyLarge?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 16),
          ModernButton(
            text: l10n.retry,
            onPressed: () {
              context.read<ProjectBloc>().add(const ProjectsLoaded());
            },
            type: ModernButtonType.outline,
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyState(AppLocalizations l10n) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.all(40),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              gradient: theme.brightness == Brightness.light
                  ? AppTheme.primaryGradient
                  : AppTheme.primaryGradientDark,
              borderRadius: BorderRadius.circular(24),
            ),
            child: const Icon(
              Icons.mic_rounded,
              size: 64,
              color: Colors.white,
            ),
          )
              .animate(onPlay: (controller) => controller.repeat(reverse: true))
              .scale(
                begin: const Offset(1.0, 1.0),
                end: const Offset(1.1, 1.1),
                duration: 2000.ms,
              ),
          const SizedBox(height: 24),
          Text(
            'Create Your First Voice Project',
            style: theme.textTheme.headlineSmall?.copyWith(
              fontWeight: FontWeight.w700,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 12),
          Text(
            'Transform text into natural-sounding speech with AI-powered voice generation',
            style: theme.textTheme.bodyLarge?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
              height: 1.5,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          GradientShimmerButton(
            text: l10n.newProject,
            onPressed: () => _showCreateProjectDialog(context),
            icon: const Icon(Icons.add, color: Colors.white),
            size: ModernButtonSize.large,
          ),
        ],
      ),
    );
  }

  int _getCrossAxisCount(BuildContext context) {
    final width = MediaQuery.of(context).size.width;
    if (width > 1200) return 4;
    if (width > 900) return 3;
    if (width > 600) return 2;
    return 1;
  }

  Color _getProjectAccentColor(int index) {
    final colors = [
      const Color(0xFF6366F1), // Indigo
      const Color(0xFF8B5CF6), // Purple
      const Color(0xFF06B6D4), // Cyan
      const Color(0xFF10B981), // Emerald
      const Color(0xFFF59E0B), // Amber
      const Color(0xFFEF4444), // Red
    ];
    return colors[index % colors.length];
  }

  void _showDeleteConfirmation(BuildContext context, Project project) {
    showDialog(
      context: context,
      builder: (context) => Dialog(
        backgroundColor: Colors.transparent,
        child: GlassmorphismContainer(
          width: 320,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.errorContainer,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Icon(
                  Icons.delete_outline,
                  size: 32,
                  color: Theme.of(context).colorScheme.onErrorContainer,
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'Delete Project',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
              ),
              const SizedBox(height: 8),
              Text(
                'Are you sure you want to delete "${project.title}"? This action cannot be undone.',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                    ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 24),
              Row(
                children: [
                  Expanded(
                    child: ModernButton(
                      text: AppLocalizations.of(context).cancel,
                      onPressed: () => Navigator.of(context).pop(),
                      type: ModernButtonType.outline,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: ModernButton(
                      text: 'Delete',
                      onPressed: () {
                        Navigator.of(context).pop();
                        context.read<ProjectBloc>().add(
                              ProjectDeleted(id: project.id),
                            );
                      },
                      type: ModernButtonType.primary,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class ProjectCard extends StatelessWidget {
  final String title;
  final String description;
  final DateTime lastModified;
  final bool hasAudio;
  final Color accentColor;
  final VoidCallback onTap;
  final VoidCallback onDelete;

  const ProjectCard({
    super.key,
    required this.title,
    required this.description,
    required this.lastModified,
    required this.hasAudio,
    required this.accentColor,
    required this.onTap,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return ModernCard(
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      title,
                      style: Theme.of(context).textTheme.titleMedium,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  PopupMenuButton<String>(
                    onSelected: (value) {
                      if (value == 'delete') {
                        onDelete();
                      }
                    },
                    itemBuilder: (context) => [
                      const PopupMenuItem(
                        value: 'delete',
                        child: Row(
                          children: [
                            Icon(Icons.delete, color: Colors.red),
                            SizedBox(width: 8),
                            Text('Delete', style: TextStyle(color: Colors.red)),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Expanded(
                child: Text(
                  description,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  if (hasAudio) const Icon(Icons.music_note, size: 16),
                  Text(
                    _formatDate(lastModified),
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _formatDate(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inDays == 0) {
      if (difference.inHours == 0) {
        return '${difference.inMinutes} min ago';
      }
      return '${difference.inHours}h ago';
    } else if (difference.inDays < 7) {
      return '${difference.inDays}d ago';
    } else {
      return '${dateTime.day}/${dateTime.month}/${dateTime.year}';
    }
  }
}
