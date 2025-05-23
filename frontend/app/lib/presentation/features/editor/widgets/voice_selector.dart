import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../../../../domain/models/voice.dart';
import '../../../blocs/editor/editor_bloc.dart';
import '../../../blocs/editor/editor_event.dart';
import '../../../blocs/editor/editor_state.dart';

class VoiceSelector extends StatelessWidget {
  const VoiceSelector({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return BlocBuilder<EditorBloc, EditorState>(
      builder: (context, state) {
        if (state.availableVoices.isEmpty) {
          return const CircularProgressIndicator();
        }

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              l10n.selectVoice,
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            SizedBox(
              height: 156,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                itemCount: state.availableVoices.length,
                itemBuilder: (context, index) {
                  final voice = state.availableVoices[index];
                  final isSelected = voice.id == state.script.voiceId;

                  return Padding(
                    padding: const EdgeInsets.only(right: 12),
                    child: VoiceCard(
                      voice: voice,
                      isSelected: isSelected,
                      onTap: () {
                        context.read<EditorBloc>().add(
                              VoiceSelected(
                                voiceId: voice.id,
                                voiceName: voice.name,
                                language: voice.language,
                              ),
                            );
                      },
                    ),
                  );
                },
              ),
            ),
          ],
        );
      },
    );
  }
}

class VoiceCard extends StatelessWidget {
  final Voice voice;
  final bool isSelected;
  final VoidCallback onTap;

  const VoiceCard({
    super.key,
    required this.voice,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        width: 120,
        decoration: BoxDecoration(
          color: isSelected
              ? theme.colorScheme.primaryContainer
              : theme.colorScheme.surfaceVariant,
          borderRadius: BorderRadius.circular(12),
          border: isSelected
              ? Border.all(color: theme.colorScheme.primary, width: 2)
              : null,
        ),
        padding: const EdgeInsets.all(8),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircleAvatar(
              radius: 32,
              backgroundImage: AssetImage(voice.avatarUrl),
            ),
            const SizedBox(height: 8),
            Text(
              voice.name,
              style: theme.textTheme.titleSmall,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 4),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  voice.type == VoiceType.cloned
                      ? Icons.person
                      : Icons.record_voice_over,
                  size: 14,
                  color: theme.colorScheme.onSurfaceVariant,
                ),
                const SizedBox(width: 4),
                Text(
                  voice.language.toUpperCase(),
                  style: theme.textTheme.bodySmall!.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
