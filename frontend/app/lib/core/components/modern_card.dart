import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:vaani/core/theme/app_theme.dart';

enum ModernCardType {
  elevated,
  outlined,
  filled,
  glassmorphism,
}

class ModernCard extends StatefulWidget {
  final Widget child;
  final VoidCallback? onTap;
  final ModernCardType type;
  final EdgeInsets? padding;
  final BorderRadius? borderRadius;
  final List<BoxShadow>? shadows;
  final Color? backgroundColor;
  final Border? border;
  final bool showShimmer;
  final Duration animationDuration;

  const ModernCard({
    super.key,
    required this.child,
    this.onTap,
    this.type = ModernCardType.elevated,
    this.padding,
    this.borderRadius,
    this.shadows,
    this.backgroundColor,
    this.border,
    this.showShimmer = false,
    this.animationDuration = const Duration(milliseconds: 300),
  });

  @override
  State<ModernCard> createState() => _ModernCardState();
}

class _ModernCardState extends State<ModernCard> with TickerProviderStateMixin {
  late AnimationController _hoverController;
  late AnimationController _tapController;
  late Animation<double> _hoverAnimation;
  late Animation<double> _tapAnimation;
  bool _isHovered = false;

  @override
  void initState() {
    super.initState();
    _hoverController = AnimationController(
      duration: widget.animationDuration,
      vsync: this,
    );
    _tapController = AnimationController(
      duration: const Duration(milliseconds: 150),
      vsync: this,
    );

    _hoverAnimation = Tween<double>(
      begin: 0,
      end: 1,
    ).animate(CurvedAnimation(
      parent: _hoverController,
      curve: Curves.easeOutCubic,
    ));

    _tapAnimation = Tween<double>(
      begin: 1,
      end: 0.98,
    ).animate(CurvedAnimation(
      parent: _tapController,
      curve: Curves.easeOut,
    ));
  }

  @override
  void dispose() {
    _hoverController.dispose();
    _tapController.dispose();
    super.dispose();
  }

  void _onHover(bool isHovered) {
    setState(() {
      _isHovered = isHovered;
    });
    if (isHovered) {
      _hoverController.forward();
    } else {
      _hoverController.reverse();
    }
  }

  void _onTapDown() {
    _tapController.forward();
  }

  void _onTapUp() {
    _tapController.reverse();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return AnimatedBuilder(
      animation: Listenable.merge([_hoverAnimation, _tapAnimation]),
      builder: (context, child) {
        return Transform.scale(
          scale: _tapAnimation.value,
          child: MouseRegion(
            onEnter: (_) => _onHover(true),
            onExit: (_) => _onHover(false),
            child: GestureDetector(
              onTapDown: (_) => _onTapDown(),
              onTapUp: (_) => _onTapUp(),
              onTapCancel: () => _onTapUp(),
              onTap: widget.onTap,
              child: AnimatedContainer(
                duration: widget.animationDuration,
                curve: Curves.easeOutCubic,
                padding: widget.padding ?? const EdgeInsets.all(20),
                decoration: _getDecoration(context, colorScheme),
                child: widget.child,
              ),
            ),
          ),
        );
      },
    );
  }

  BoxDecoration _getDecoration(BuildContext context, ColorScheme colorScheme) {
    final theme = Theme.of(context);
    final hoverValue = _hoverAnimation.value;

    switch (widget.type) {
      case ModernCardType.elevated:
        return BoxDecoration(
          color: widget.backgroundColor ?? colorScheme.surface,
          borderRadius: widget.borderRadius ?? BorderRadius.circular(20),
          border: widget.border,
          boxShadow: widget.shadows ??
              [
                BoxShadow(
                  color: colorScheme.shadow
                      .withOpacity(0.08 + (hoverValue * 0.12)),
                  offset: Offset(0, 4 + (hoverValue * 8)),
                  blurRadius: 16 + (hoverValue * 16),
                  spreadRadius: 0,
                ),
                BoxShadow(
                  color: colorScheme.primary
                      .withOpacity(0.05 + (hoverValue * 0.1)),
                  offset: Offset(0, 1 + (hoverValue * 2)),
                  blurRadius: 6 + (hoverValue * 6),
                  spreadRadius: 0,
                ),
              ],
        );

      case ModernCardType.outlined:
        return BoxDecoration(
          color: widget.backgroundColor ?? colorScheme.surface,
          borderRadius: widget.borderRadius ?? BorderRadius.circular(20),
          border: widget.border ??
              Border.all(
                color:
                    colorScheme.outline.withOpacity(0.5 + (hoverValue * 0.5)),
                width: 1 + (hoverValue * 0.5),
              ),
          boxShadow: hoverValue > 0
              ? [
                  BoxShadow(
                    color: colorScheme.primary.withOpacity(hoverValue * 0.1),
                    offset: Offset(0, 2 * hoverValue),
                    blurRadius: 8 * hoverValue,
                    spreadRadius: 0,
                  ),
                ]
              : null,
        );

      case ModernCardType.filled:
        return BoxDecoration(
          color: widget.backgroundColor ?? colorScheme.surfaceVariant,
          borderRadius: widget.borderRadius ?? BorderRadius.circular(20),
          border: widget.border,
          boxShadow: hoverValue > 0
              ? [
                  BoxShadow(
                    color: colorScheme.shadow.withOpacity(hoverValue * 0.1),
                    offset: Offset(0, 4 * hoverValue),
                    blurRadius: 12 * hoverValue,
                    spreadRadius: 0,
                  ),
                ]
              : null,
        );

      case ModernCardType.glassmorphism:
        return BoxDecoration(
          color: (widget.backgroundColor ?? colorScheme.surface)
              .withOpacity(0.7 + (hoverValue * 0.1)),
          borderRadius: widget.borderRadius ?? BorderRadius.circular(20),
          border: widget.border ??
              Border.all(
                color: theme.brightness == Brightness.light
                    ? Colors.white.withOpacity(0.2 + (hoverValue * 0.3))
                    : Colors.white.withOpacity(0.1 + (hoverValue * 0.2)),
                width: 1,
              ),
          boxShadow: [
            BoxShadow(
              color: theme.brightness == Brightness.light
                  ? Colors.black.withOpacity(0.05 + (hoverValue * 0.1))
                  : Colors.black.withOpacity(0.2 + (hoverValue * 0.2)),
              offset: Offset(0, 8 + (hoverValue * 8)),
              blurRadius: 32 + (hoverValue * 16),
              spreadRadius: 0,
            ),
          ],
        );
    }
  }
}

// Project card specifically for home screen
class ProjectCard extends StatelessWidget {
  final String title;
  final String description;
  final DateTime lastModified;
  final VoidCallback? onTap;
  final VoidCallback? onDelete;
  final bool hasAudio;
  final Color? accentColor;

  const ProjectCard({
    super.key,
    required this.title,
    required this.description,
    required this.lastModified,
    this.onTap,
    this.onDelete,
    this.hasAudio = false,
    this.accentColor,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return ModernCard(
      type: ModernCardType.elevated,
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 4,
                height: 32,
                decoration: BoxDecoration(
                  color: accentColor ?? colorScheme.primary,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  title,
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              PopupMenuButton<String>(
                icon: Icon(
                  Icons.more_vert,
                  color: colorScheme.onSurfaceVariant,
                ),
                onSelected: (value) {
                  if (value == 'delete' && onDelete != null) {
                    onDelete!();
                  }
                },
                itemBuilder: (context) => [
                  PopupMenuItem(
                    value: 'delete',
                    child: Row(
                      children: [
                        Icon(Icons.delete_outline, color: colorScheme.error),
                        const SizedBox(width: 12),
                        Text(
                          'Delete',
                          style: TextStyle(color: colorScheme.error),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ],
          ),
          const SizedBox(height: 12),
          Expanded(
            child: Text(
              description,
              style: theme.textTheme.bodyMedium?.copyWith(
                color: colorScheme.onSurfaceVariant,
                height: 1.5,
              ),
              maxLines: 3,
              overflow: TextOverflow.ellipsis,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              if (hasAudio)
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: colorScheme.primaryContainer,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        Icons.music_note,
                        size: 14,
                        color: colorScheme.onPrimaryContainer,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        'Audio',
                        style: theme.textTheme.labelSmall?.copyWith(
                          color: colorScheme.onPrimaryContainer,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ],
                  ),
                ),
              if (!hasAudio) const Spacer(),
              Text(
                _formatDate(lastModified),
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colorScheme.onSurfaceVariant.withOpacity(0.7),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inDays == 0) {
      if (difference.inHours == 0) {
        return '${difference.inMinutes}m ago';
      }
      return '${difference.inHours}h ago';
    } else if (difference.inDays < 7) {
      return '${difference.inDays}d ago';
    } else {
      return '${dateTime.day}/${dateTime.month}/${dateTime.year}';
    }
  }
}

// Stats card for dashboard
class StatsCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color? iconColor;
  final Widget? trend;

  const StatsCard({
    super.key,
    required this.title,
    required this.value,
    required this.icon,
    this.iconColor,
    this.trend,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return ModernCard(
      type: ModernCardType.filled,
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: (iconColor ?? colorScheme.primary).withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  icon,
                  color: iconColor ?? colorScheme.primary,
                  size: 24,
                ),
              ),
              const Spacer(),
              if (trend != null) trend!,
            ],
          ),
          const SizedBox(height: 16),
          Text(
            value,
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.w700,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            title,
            style: theme.textTheme.bodyMedium?.copyWith(
              color: colorScheme.onSurfaceVariant,
            ),
          ),
        ],
      ),
    );
  }
}

// Glassmorphism container for special effects
class GlassmorphismContainer extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsets? padding;
  final BorderRadius? borderRadius;
  final double opacity;

  const GlassmorphismContainer({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.padding,
    this.borderRadius,
    this.opacity = 0.7,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Container(
      width: width,
      height: height,
      padding: padding ?? const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: colorScheme.surface.withOpacity(opacity),
        borderRadius: borderRadius ?? BorderRadius.circular(20),
        border: Border.all(
          color: theme.brightness == Brightness.light
              ? Colors.white.withOpacity(0.2)
              : Colors.white.withOpacity(0.1),
          width: 1,
        ),
        boxShadow: [
          BoxShadow(
            color: theme.brightness == Brightness.light
                ? Colors.black.withOpacity(0.1)
                : Colors.black.withOpacity(0.3),
            offset: const Offset(0, 8),
            blurRadius: 32,
            spreadRadius: 0,
          ),
        ],
      ),
      child: child,
    );
  }
}
