import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../../../blocs/editor/editor_bloc.dart';
import '../../../blocs/editor/editor_event.dart';
import '../../../blocs/editor/editor_state.dart';

class RichTextField extends StatefulWidget {
  const RichTextField({super.key});

  @override
  State<RichTextField> createState() => _RichTextFieldState();
}

class _RichTextFieldState extends State<RichTextField> {
  late TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return BlocConsumer<EditorBloc, EditorState>(
      listenWhen: (previous, current) =>
          previous.script.text != current.script.text,
      listener: (context, state) {
        // Update the controller if the text changed from elsewhere (e.g. auto-punctuation)
        if (_controller.text != state.script.text) {
          _controller.text = state.script.text;
          // Place cursor at the end
          _controller.selection = TextSelection.fromPosition(
            TextPosition(offset: _controller.text.length),
          );
        }
      },
      builder: (context, state) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _controller,
              decoration: InputDecoration(
                hintText: l10n.editScript,
                border: const OutlineInputBorder(),
                contentPadding: const EdgeInsets.all(16),
                suffixIcon: IconButton(
                  icon: const Icon(Icons.settings),
                  tooltip: l10n.settings,
                  onPressed: () {
                    context.read<EditorBloc>().add(ApplyAutoPunctuation());
                  },
                ),
              ),
              maxLines: 8,
              onChanged: (value) {
                context.read<EditorBloc>().add(TextChanged(value));
              },
              style: theme.textTheme.bodyLarge,
            ),
            const SizedBox(height: 8),
            // Character count
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Text(
                  l10n.characterCount(state.script.charCount),
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: state.isOverCharLimit
                        ? theme.colorScheme.error
                        : theme.colorScheme.onSurface,
                  ),
                ),
              ],
            ),
          ],
        );
      },
    );
  }
}
