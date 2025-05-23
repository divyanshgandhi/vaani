import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:vaani/core/components/modern_button.dart';
import 'package:vaani/core/components/modern_card.dart';
import 'package:vaani/core/theme/app_theme.dart';
import 'package:vaani/core/animations/app_animations.dart';
import '../../../data/repositories/pricing_repository_impl.dart';
import '../../../data/repositories/voice_repository_impl.dart';
import '../../../domain/repositories/pricing_repository.dart';
import '../../../domain/repositories/voice_repository.dart';
import '../../blocs/editor/editor_bloc.dart';
import '../../blocs/editor/editor_state.dart';
import 'widgets/emotion_slider.dart';
import 'widgets/rich_text_field.dart';
import 'widgets/voice_selector.dart';

class EditorScreen extends StatefulWidget {
  final String? projectId;

  const EditorScreen({super.key, this.projectId});

  @override
  State<EditorScreen> createState() => _EditorScreenState();
}

class _EditorScreenState extends State<EditorScreen>
    with TickerProviderStateMixin {
  late AnimationController _headerAnimationController;
  late AnimationController _contentAnimationController;
  late AnimationController _fabAnimationController;

  @override
  void initState() {
    super.initState();
    _headerAnimationController = AnimationController(
      duration: const Duration(milliseconds: 600),
      vsync: this,
    );
    _contentAnimationController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    _fabAnimationController = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );

    // Start animations
    _headerAnimationController.forward();
    Future.delayed(const Duration(milliseconds: 200), () {
      if (mounted) _contentAnimationController.forward();
    });
    Future.delayed(const Duration(milliseconds: 600), () {
      if (mounted) _fabAnimationController.forward();
    });
  }

  @override
  void dispose() {
    _headerAnimationController.dispose();
    _contentAnimationController.dispose();
    _fabAnimationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final dio = Dio(); // In a real app, this would be injected

    return MultiRepositoryProvider(
      providers: [
        RepositoryProvider<VoiceRepository>(
          create: (context) => VoiceRepositoryImpl(dio),
        ),
        RepositoryProvider<PricingRepository>(
          create: (context) => PricingRepositoryImpl(dio),
        ),
      ],
      child: BlocProvider(
        create: (context) => EditorBloc(
          voiceRepository: context.read<VoiceRepository>(),
          pricingRepository: context.read<PricingRepository>(),
        ),
        child: BlocBuilder<EditorBloc, EditorState>(
          builder: (context, state) {
            return Scaffold(
              backgroundColor: Theme.of(context).colorScheme.background,
              body: CustomScrollView(
                slivers: [
                  // Modern App Bar
                  _buildModernAppBar(context, l10n, state),

                  // Editor Content
                  SliverToBoxAdapter(
                    child: state.status == EditorStatus.initial ||
                            state.status == EditorStatus.loading
                        ? _buildLoadingState(context)
                        : _buildEditorContent(context, state, l10n),
                  ),
                ],
              ),
              floatingActionButton:
                  _buildFloatingActionButton(context, state, l10n),
              floatingActionButtonLocation:
                  FloatingActionButtonLocation.centerFloat,
            );
          },
        ),
      ),
    );
  }

  Widget _buildModernAppBar(
      BuildContext context, AppLocalizations l10n, EditorState state) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return SliverAppBar(
      expandedHeight: 140,
      floating: true,
      pinned: true,
      backgroundColor: Colors.transparent,
      elevation: 0,
      leading: Container(
        margin: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.2),
          borderRadius: BorderRadius.circular(12),
        ),
        child: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      actions: [
        Container(
          margin: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.2),
            borderRadius: BorderRadius.circular(12),
          ),
          child: IconButton(
            icon: const Icon(Icons.save, color: Colors.white),
            onPressed:
                state.status == EditorStatus.ready && !state.isOverCharLimit
                    ? () {
                        // TODO: Save project
                        Navigator.of(context).pop();
                      }
                    : null,
          ),
        ),
      ],
      flexibleSpace: FlexibleSpaceBar(
        background: Container(
          decoration: BoxDecoration(
            gradient: theme.brightness == Brightness.light
                ? AppTheme.primaryGradient
                : AppTheme.primaryGradientDark,
          ),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(20, 60, 20, 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Text(
                    widget.projectId == null ? 'New Project' : 'Create Audio',
                    style: theme.textTheme.headlineMedium?.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.w700,
                    ),
                  )
                      .animate(controller: _headerAnimationController)
                      .fadeIn(duration: 600.ms)
                      .slideY(begin: 0.3, end: 0),
                  const SizedBox(height: 8),
                  Text(
                    'Transform your text into natural-sounding audio',
                    style: theme.textTheme.bodyLarge?.copyWith(
                      color: Colors.white.withOpacity(0.9),
                    ),
                  )
                      .animate(controller: _headerAnimationController)
                      .fadeIn(duration: 600.ms, delay: 200.ms)
                      .slideY(begin: 0.3, end: 0),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildLoadingState(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      height: MediaQuery.of(context).size.height * 0.6,
      padding: const EdgeInsets.all(40),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
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
              size: 48,
              color: Colors.white,
            ),
          ).animate(onPlay: (controller) => controller.repeat()).scale(
                begin: const Offset(1.0, 1.0),
                end: const Offset(1.1, 1.1),
                duration: 1500.ms,
              ),
          const SizedBox(height: 24),
          Text(
            'Preparing your voice studio...',
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          AppAnimations.loadingDots(
            color: theme.colorScheme.primary,
          ),
        ],
      ),
    );
  }

  Widget _buildEditorContent(
      BuildContext context, EditorState state, AppLocalizations l10n) {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Pricing Banner
          _buildPricingBanner(context, state)
              .animate(controller: _contentAnimationController)
              .fadeIn(duration: 400.ms)
              .slideY(begin: 0.3, end: 0),

          const SizedBox(height: 24),

          // Text Editor Section
          _buildTextEditorSection(context)
              .animate(controller: _contentAnimationController)
              .fadeIn(duration: 400.ms, delay: 100.ms)
              .slideY(begin: 0.3, end: 0),

          const SizedBox(height: 24),

          // Voice & Emotion Controls
          _buildVoiceEmotionSection(context)
              .animate(controller: _contentAnimationController)
              .fadeIn(duration: 400.ms, delay: 200.ms)
              .slideY(begin: 0.3, end: 0),

          const SizedBox(height: 100), // Space for FAB
        ],
      ),
    );
  }

  Widget _buildPricingBanner(BuildContext context, EditorState state) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return ModernCard(
      type: ModernCardType.glassmorphism,
      backgroundColor: colorScheme.primaryContainer,
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: colorScheme.primary.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(
              Icons.monetization_on,
              color: colorScheme.primary,
              size: 24,
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Cost Estimate',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: colorScheme.onPrimaryContainer,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '${state.script.charCount} characters × ₹${state.pricePerChar.toStringAsFixed(4)}',
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: colorScheme.onPrimaryContainer.withOpacity(0.8),
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(
              color: colorScheme.primary,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Text(
              '₹${state.totalCost.toStringAsFixed(2)}',
              style: theme.textTheme.titleMedium?.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.w700,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTextEditorSection(BuildContext context) {
    final theme = Theme.of(context);

    return ModernCard(
      type: ModernCardType.elevated,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                Icons.edit_rounded,
                color: theme.colorScheme.primary,
                size: 24,
              ),
              const SizedBox(width: 12),
              Text(
                'Text to Audio',
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const Spacer(),
              ModernButton(
                text: 'Auto-Punctuation',
                type: ModernButtonType.ghost,
                size: ModernButtonSize.small,
                icon: const Icon(Icons.auto_fix_high, size: 16),
                onPressed: () {
                  // TODO: Implement auto-punctuation
                },
              ),
            ],
          ),
          const SizedBox(height: 16),
          const RichTextField(),
        ],
      ),
    );
  }

  Widget _buildVoiceEmotionSection(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Voice Selector
        Expanded(
          flex: 3,
          child: ModernCard(
            type: ModernCardType.elevated,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      Icons.record_voice_over,
                      color: Theme.of(context).colorScheme.secondary,
                      size: 24,
                    ),
                    const SizedBox(width: 12),
                    Text(
                      'Voice Selection',
                      style: Theme.of(context).textTheme.titleLarge?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                const VoiceSelector(),
              ],
            ),
          ),
        ),

        const SizedBox(width: 16),

        // Emotion Slider
        Expanded(
          flex: 2,
          child: ModernCard(
            type: ModernCardType.elevated,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      Icons.sentiment_satisfied,
                      color: Theme.of(context).colorScheme.tertiary,
                      size: 24,
                    ),
                    const SizedBox(width: 12),
                    Text(
                      'Emotion',
                      style: Theme.of(context).textTheme.titleLarge?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                const EmotionSlider(),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildFloatingActionButton(
      BuildContext context, EditorState state, AppLocalizations l10n) {
    return AnimatedBuilder(
      animation: _fabAnimationController,
      builder: (context, child) {
        return Transform.scale(
          scale: _fabAnimationController.value,
          child: Container(
            margin: const EdgeInsets.symmetric(horizontal: 20),
            child: SizedBox(
              width: double.infinity,
              child: GradientShimmerButton(
                text: l10n.playPreview,
                icon: const Icon(Icons.play_arrow, color: Colors.white),
                size: ModernButtonSize.large,
                onPressed: state.status == EditorStatus.ready &&
                        state.script.text.isNotEmpty &&
                        !state.isOverCharLimit
                    ? () {
                        // TODO: Implement preview functionality (Epic FE-5)
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(
                            content: const Text(
                                'Preview functionality coming soon!'),
                            backgroundColor:
                                Theme.of(context).colorScheme.primaryContainer,
                          ),
                        );
                      }
                    : null,
              ),
            ),
          ),
        );
      },
    );
  }
}
