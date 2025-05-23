import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:vaani/presentation/blocs/tts/tts_bloc.dart';
import 'package:vaani/presentation/blocs/tts/tts_event.dart';
import 'package:vaani/presentation/blocs/tts/tts_state.dart';
import 'package:just_audio/just_audio.dart';

class TtsDemoScreen extends StatefulWidget {
  const TtsDemoScreen({super.key});

  @override
  State<TtsDemoScreen> createState() => _TtsDemoScreenState();
}

class _TtsDemoScreenState extends State<TtsDemoScreen> {
  final TextEditingController _textController = TextEditingController();
  final AudioPlayer _audioPlayer = AudioPlayer();

  String _selectedVoiceId = 'hi_female_1';
  int _emotion = 50;
  double _speed = 1.0;
  double _pitch = 1.0;

  @override
  void initState() {
    super.initState();
    _textController.text = 'नमस्ते दुनिया! यह एक टेस्ट मैसेज है।';
    // Load voices when screen initializes
    context.read<TtsBloc>().add(const TtsVoicesRequested());
  }

  @override
  void dispose() {
    _textController.dispose();
    _audioPlayer.dispose();
    super.dispose();
  }

  void _detectLanguage() {
    if (_textController.text.isNotEmpty) {
      context
          .read<TtsBloc>()
          .add(TtsLanguageDetectionRequested(_textController.text));
    }
  }

  void _generateTts() {
    if (_textController.text.isNotEmpty) {
      context.read<TtsBloc>().add(TtsGenerationRequested(
            text: _textController.text,
            voiceId: _selectedVoiceId,
            emotion: _emotion,
            speed: _speed,
            pitch: _pitch,
          ));
    }
  }

  void _previewTts() {
    if (_textController.text.isNotEmpty) {
      context.read<TtsBloc>().add(TtsPreviewRequested(
            text: _textController.text,
            voiceId: _selectedVoiceId,
            emotion: _emotion,
            speed: _speed,
            pitch: _pitch,
          ));
    }
  }

  void _playAudio(String url) async {
    try {
      await _audioPlayer.setUrl(url);
      await _audioPlayer.play();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing audio: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => GetIt.instance<TtsBloc>(),
      child: Scaffold(
        appBar: AppBar(
          title: const Text('TTS Demo'),
          backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        ),
        body: Padding(
          padding: const EdgeInsets.all(16.0),
          child: BlocConsumer<TtsBloc, TtsState>(
            listener: (context, state) {
              if (state is TtsGenerationFailure ||
                  state is TtsLanguageDetectionFailure ||
                  state is TtsVoicesFailure ||
                  state is TtsPreviewFailure) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Error: ${_getErrorMessage(state)}')),
                );
              }

              if (state is TtsPreviewCompleted) {
                _playAudio(state.outputUrl);
              }
            },
            builder: (context, state) {
              return Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // Text Input
                  TextField(
                    controller: _textController,
                    decoration: const InputDecoration(
                      labelText: 'Enter text to convert to speech',
                      border: OutlineInputBorder(),
                    ),
                    maxLines: 3,
                  ),
                  const SizedBox(height: 16),

                  // Voice Selection
                  if (state is TtsVoicesSuccess) ...[
                    DropdownButtonFormField<String>(
                      value: _selectedVoiceId,
                      decoration: const InputDecoration(
                        labelText: 'Select Voice',
                        border: OutlineInputBorder(),
                      ),
                      items: state.voices.map((voice) {
                        return DropdownMenuItem(
                          value: voice.id,
                          child: Text('${voice.name} (${voice.language})'),
                        );
                      }).toList(),
                      onChanged: (value) {
                        if (value != null) {
                          setState(() {
                            _selectedVoiceId = value;
                          });
                        }
                      },
                    ),
                    const SizedBox(height: 16),
                  ],

                  // Emotion Slider
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('Emotion: $_emotion'),
                      Slider(
                        value: _emotion.toDouble(),
                        min: 0,
                        max: 100,
                        divisions: 100,
                        onChanged: (value) {
                          setState(() {
                            _emotion = value.toInt();
                          });
                        },
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // Speed Slider
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('Speed: ${_speed.toStringAsFixed(1)}x'),
                      Slider(
                        value: _speed,
                        min: 0.5,
                        max: 2.0,
                        divisions: 15,
                        onChanged: (value) {
                          setState(() {
                            _speed = value;
                          });
                        },
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // Action Buttons
                  Wrap(
                    spacing: 8,
                    children: [
                      ElevatedButton(
                        onPressed: _detectLanguage,
                        child: const Text('Detect Language'),
                      ),
                      ElevatedButton(
                        onPressed: _previewTts,
                        child: const Text('Preview'),
                      ),
                      ElevatedButton(
                        onPressed: _generateTts,
                        child: const Text('Generate'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // Status Display
                  Expanded(
                    child: Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.grey[100],
                        border: Border.all(color: Colors.grey),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: _buildStatusWidget(state),
                    ),
                  ),
                ],
              );
            },
          ),
        ),
      ),
    );
  }

  Widget _buildStatusWidget(TtsState state) {
    if (state is TtsInitial) {
      return const Text('Ready to generate speech');
    } else if (state is TtsLanguageDetectionLoading) {
      return const Column(
        children: [
          CircularProgressIndicator(),
          SizedBox(height: 8),
          Text('Detecting language...'),
        ],
      );
    } else if (state is TtsLanguageDetectionSuccess) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Language: ${state.languageDetection.language}'),
          Text('Is Indic: ${state.languageDetection.isIndic}'),
        ],
      );
    } else if (state is TtsGenerationLoading) {
      return const Column(
        children: [
          CircularProgressIndicator(),
          SizedBox(height: 8),
          Text('Generating speech...'),
        ],
      );
    } else if (state is TtsGenerationSuccess) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Job ID: ${state.job.id}'),
          Text('Status: ${state.job.status}'),
          Text('Created: ${state.job.createdAt}'),
          if (state.job.outputUrl != null) ...[
            const SizedBox(height: 8),
            ElevatedButton(
              onPressed: () => _playAudio(state.job.outputUrl!),
              child: const Text('Play Audio'),
            ),
          ],
        ],
      );
    } else if (state is TtsPreviewStarted) {
      return const Column(
        children: [
          CircularProgressIndicator(),
          SizedBox(height: 8),
          Text('Starting preview...'),
        ],
      );
    } else if (state is TtsPreviewProgress) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const CircularProgressIndicator(),
          const SizedBox(height: 8),
          Text('Preview in progress...'),
          Text('Status: ${state.data['status']}'),
          if (state.data['progress'] != null)
            Text('Progress: ${state.data['progress']}%'),
        ],
      );
    } else if (state is TtsPreviewCompleted) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Preview completed!'),
          Text('Job ID: ${state.jobId}'),
          if (state.duration != null)
            Text('Duration: ${state.duration!.toStringAsFixed(1)}s'),
          const SizedBox(height: 8),
          ElevatedButton(
            onPressed: () => _playAudio(state.outputUrl),
            child: const Text('Play Audio'),
          ),
        ],
      );
    } else if (state is TtsVoicesLoading) {
      return const Column(
        children: [
          CircularProgressIndicator(),
          SizedBox(height: 8),
          Text('Loading voices...'),
        ],
      );
    } else {
      return Text('Current state: ${state.runtimeType}');
    }
  }

  String _getErrorMessage(TtsState state) {
    if (state is TtsGenerationFailure) return state.error;
    if (state is TtsLanguageDetectionFailure) return state.error;
    if (state is TtsVoicesFailure) return state.error;
    if (state is TtsPreviewFailure) return state.error;
    return 'Unknown error';
  }
}
