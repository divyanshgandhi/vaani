import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../../../data/repositories/pricing_repository_impl.dart';
import '../../../data/repositories/voice_repository_impl.dart';
import '../../../domain/repositories/pricing_repository.dart';
import '../../../domain/repositories/voice_repository.dart';
import '../../blocs/editor/editor_bloc.dart';
import '../../blocs/editor/editor_state.dart';
import 'widgets/emotion_slider.dart';
import 'widgets/pricing_info.dart';
import 'widgets/rich_text_field.dart';
import 'widgets/voice_selector.dart';

class EditorScreen extends StatelessWidget {
  final String? projectId;

  const EditorScreen({super.key, this.projectId});

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
              appBar: AppBar(
                title:
                    Text(projectId == null ? l10n.newProject : l10n.editScript),
                leading: IconButton(
                  icon: const Icon(Icons.arrow_back),
                  onPressed: () => Navigator.of(context).pop(),
                ),
                actions: [
                  TextButton.icon(
                    icon: const Icon(Icons.check),
                    label: Text(l10n.save),
                    onPressed: state.status == EditorStatus.ready &&
                            !state.isOverCharLimit
                        ? () {
                            // TODO: Save project (will be implemented in a separate story)
                            Navigator.of(context).pop();
                          }
                        : null,
                  ),
                ],
              ),
              body: state.status == EditorStatus.initial ||
                      state.status == EditorStatus.loading
                  ? const Center(child: CircularProgressIndicator())
                  : _buildEditorContent(context, state),
            );
          },
        ),
      ),
    );
  }

  Widget _buildEditorContent(BuildContext context, EditorState state) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Text editor
          const RichTextField(),
          const SizedBox(height: 24),

          // Pricing information
          const PricingInfo(),
          const SizedBox(height: 24),

          // Voice selector
          const VoiceSelector(),
          const SizedBox(height: 24),

          // Emotion slider
          const EmotionSlider(),
          const SizedBox(height: 32),

          // Preview button (disabled for now)
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              icon: const Icon(Icons.play_arrow),
              label: Text(AppLocalizations.of(context).playPreview),
              onPressed: null, // Will be implemented in a later story
            ),
          ),
        ],
      ),
    );
  }
}
