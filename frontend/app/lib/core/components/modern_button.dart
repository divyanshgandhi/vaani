import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:vaani/core/theme/app_theme.dart';
import 'package:vaani/core/animations/app_animations.dart';

enum ModernButtonType {
  primary,
  secondary,
  outline,
  ghost,
  gradient,
}

enum ModernButtonSize {
  small,
  medium,
  large,
}

class ModernButton extends StatefulWidget {
  final String text;
  final VoidCallback? onPressed;
  final ModernButtonType type;
  final ModernButtonSize size;
  final Widget? icon;
  final bool isLoading;
  final bool isEnabled;
  final BorderRadius? borderRadius;
  final List<BoxShadow>? shadows;

  const ModernButton({
    super.key,
    required this.text,
    this.onPressed,
    this.type = ModernButtonType.primary,
    this.size = ModernButtonSize.medium,
    this.icon,
    this.isLoading = false,
    this.isEnabled = true,
    this.borderRadius,
    this.shadows,
  });

  @override
  State<ModernButton> createState() => _ModernButtonState();
}

class _ModernButtonState extends State<ModernButton>
    with TickerProviderStateMixin {
  late AnimationController _scaleController;
  late AnimationController _shimmerController;
  late Animation<double> _scaleAnimation;

  @override
  void initState() {
    super.initState();
    _scaleController = AnimationController(
      duration: const Duration(milliseconds: 150),
      vsync: this,
    );
    _shimmerController = AnimationController(
      duration: const Duration(milliseconds: 1500),
      vsync: this,
    );
    _scaleAnimation = Tween<double>(
      begin: 1.0,
      end: 0.95,
    ).animate(CurvedAnimation(
      parent: _scaleController,
      curve: Curves.easeOut,
    ));
  }

  @override
  void dispose() {
    _scaleController.dispose();
    _shimmerController.dispose();
    super.dispose();
  }

  void _onTapDown(TapDownDetails details) {
    if (widget.isEnabled && !widget.isLoading) {
      _scaleController.forward();
    }
  }

  void _onTapUp(TapUpDetails details) {
    if (widget.isEnabled && !widget.isLoading) {
      _scaleController.reverse();
    }
  }

  void _onTapCancel() {
    if (widget.isEnabled && !widget.isLoading) {
      _scaleController.reverse();
    }
  }

  EdgeInsets get _padding {
    switch (widget.size) {
      case ModernButtonSize.small:
        return const EdgeInsets.symmetric(horizontal: 16, vertical: 8);
      case ModernButtonSize.medium:
        return const EdgeInsets.symmetric(horizontal: 24, vertical: 12);
      case ModernButtonSize.large:
        return const EdgeInsets.symmetric(horizontal: 32, vertical: 16);
    }
  }

  double get _fontSize {
    switch (widget.size) {
      case ModernButtonSize.small:
        return 14;
      case ModernButtonSize.medium:
        return 16;
      case ModernButtonSize.large:
        return 18;
    }
  }

  BorderRadius get _borderRadius {
    return widget.borderRadius ??
        BorderRadius.circular(widget.size == ModernButtonSize.large ? 16 : 12);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return AnimatedBuilder(
      animation: _scaleAnimation,
      builder: (context, child) {
        return Transform.scale(
          scale: _scaleAnimation.value,
          child: GestureDetector(
            onTapDown: _onTapDown,
            onTapUp: _onTapUp,
            onTapCancel: _onTapCancel,
            onTap:
                widget.isEnabled && !widget.isLoading ? widget.onPressed : null,
            child: Container(
              padding: _padding,
              decoration: _getDecoration(context, colorScheme),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (widget.icon != null && !widget.isLoading) ...[
                    widget.icon!,
                    const SizedBox(width: 8),
                  ],
                  if (widget.isLoading) ...[
                    SizedBox(
                      width: _fontSize,
                      height: _fontSize,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        valueColor: AlwaysStoppedAnimation<Color>(
                          _getTextColor(colorScheme),
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                  ],
                  Text(
                    widget.text,
                    style: TextStyle(
                      fontSize: _fontSize,
                      fontWeight: FontWeight.w600,
                      color: _getTextColor(colorScheme),
                      letterSpacing: 0.5,
                      fontFamily: 'Inter',
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  BoxDecoration _getDecoration(BuildContext context, ColorScheme colorScheme) {
    final theme = Theme.of(context);

    switch (widget.type) {
      case ModernButtonType.primary:
        return BoxDecoration(
          color: widget.isEnabled ? colorScheme.primary : colorScheme.outline,
          borderRadius: _borderRadius,
          boxShadow: widget.isEnabled
              ? widget.shadows ??
                  [
                    BoxShadow(
                      color: colorScheme.primary.withOpacity(0.3),
                      offset: const Offset(0, 4),
                      blurRadius: 12,
                      spreadRadius: 0,
                    ),
                  ]
              : null,
        );

      case ModernButtonType.secondary:
        return BoxDecoration(
          color: widget.isEnabled ? colorScheme.secondary : colorScheme.outline,
          borderRadius: _borderRadius,
          boxShadow: widget.isEnabled
              ? widget.shadows ??
                  [
                    BoxShadow(
                      color: colorScheme.secondary.withOpacity(0.3),
                      offset: const Offset(0, 4),
                      blurRadius: 12,
                      spreadRadius: 0,
                    ),
                  ]
              : null,
        );

      case ModernButtonType.outline:
        return BoxDecoration(
          color: Colors.transparent,
          borderRadius: _borderRadius,
          border: Border.all(
            color: widget.isEnabled ? colorScheme.primary : colorScheme.outline,
            width: 1.5,
          ),
        );

      case ModernButtonType.ghost:
        return BoxDecoration(
          color: widget.isEnabled
              ? colorScheme.primary.withOpacity(0.1)
              : colorScheme.outline.withOpacity(0.1),
          borderRadius: _borderRadius,
        );

      case ModernButtonType.gradient:
        return BoxDecoration(
          gradient: widget.isEnabled
              ? (theme.brightness == Brightness.light
                  ? AppTheme.primaryGradient
                  : AppTheme.primaryGradientDark)
              : null,
          color: !widget.isEnabled ? colorScheme.outline : null,
          borderRadius: _borderRadius,
          boxShadow: widget.isEnabled
              ? widget.shadows ??
                  [
                    BoxShadow(
                      color: colorScheme.primary.withOpacity(0.4),
                      offset: const Offset(0, 8),
                      blurRadius: 24,
                      spreadRadius: 0,
                    ),
                  ]
              : null,
        );
    }
  }

  Color _getTextColor(ColorScheme colorScheme) {
    if (!widget.isEnabled) {
      return colorScheme.onSurface.withOpacity(0.38);
    }

    switch (widget.type) {
      case ModernButtonType.primary:
      case ModernButtonType.secondary:
      case ModernButtonType.gradient:
        return colorScheme.onPrimary;
      case ModernButtonType.outline:
      case ModernButtonType.ghost:
        return colorScheme.primary;
    }
  }
}

// Specialized gradient button with shimmer effect
class GradientShimmerButton extends StatelessWidget {
  final String text;
  final VoidCallback? onPressed;
  final Widget? icon;
  final bool isLoading;
  final ModernButtonSize size;

  const GradientShimmerButton({
    super.key,
    required this.text,
    this.onPressed,
    this.icon,
    this.isLoading = false,
    this.size = ModernButtonSize.medium,
  });

  @override
  Widget build(BuildContext context) {
    return ModernButton(
      text: text,
      onPressed: onPressed,
      type: ModernButtonType.gradient,
      size: size,
      icon: icon,
      isLoading: isLoading,
    ).animate(onPlay: (controller) => controller.repeat(reverse: true)).shimmer(
          duration: const Duration(milliseconds: 2000),
          color: Colors.white.withOpacity(0.3),
        );
  }
}

// Floating action button with modern styling
class ModernFloatingActionButton extends StatelessWidget {
  final VoidCallback? onPressed;
  final Widget child;
  final bool mini;
  final String? heroTag;

  const ModernFloatingActionButton({
    super.key,
    this.onPressed,
    required this.child,
    this.mini = false,
    this.heroTag,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final size = mini ? 48.0 : 56.0;

    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        gradient: theme.brightness == Brightness.light
            ? AppTheme.primaryGradient
            : AppTheme.primaryGradientDark,
        borderRadius: BorderRadius.circular(size / 2),
        boxShadow: [
          BoxShadow(
            color: theme.colorScheme.primary.withOpacity(0.4),
            offset: const Offset(0, 8),
            blurRadius: 24,
            spreadRadius: 0,
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(size / 2),
          onTap: onPressed,
          child: Center(child: child),
        ),
      ),
    )
        .animate()
        .scale(
          begin: const Offset(0.8, 0.8),
          duration: const Duration(milliseconds: 300),
          curve: Curves.elasticOut,
        )
        .fadeIn(duration: const Duration(milliseconds: 200));
  }
}
