import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';

class AppAnimations {
  // Animation durations
  static const Duration fast = Duration(milliseconds: 200);
  static const Duration medium = Duration(milliseconds: 300);
  static const Duration slow = Duration(milliseconds: 500);
  static const Duration verySlow = Duration(milliseconds: 800);

  // Easing curves
  static const Curve easeInOutCubic = Cubic(0.645, 0.045, 0.355, 1.0);
  static const Curve easeOutQuart = Cubic(0.165, 0.84, 0.44, 1.0);
  static const Curve easeInOutQuart = Cubic(0.77, 0, 0.175, 1.0);
  static const Curve bounceOut = Curves.bounceOut;

  // Common animations
  static List<Effect> get fadeInUp => [
        FadeEffect(
          begin: 0,
          end: 1,
          duration: medium,
          curve: easeOutQuart,
        ),
        SlideEffect(
          begin: const Offset(0, 0.3),
          end: Offset.zero,
          duration: medium,
          curve: easeOutQuart,
        ),
      ];

  static List<Effect> get fadeInScale => [
        FadeEffect(
          begin: 0,
          end: 1,
          duration: medium,
          curve: easeOutQuart,
        ),
        ScaleEffect(
          begin: const Offset(0.8, 0.8),
          end: const Offset(1.0, 1.0),
          duration: medium,
          curve: easeOutQuart,
        ),
      ];

  static List<Effect> get slideInFromRight => [
        SlideEffect(
          begin: const Offset(1.0, 0),
          end: Offset.zero,
          duration: medium,
          curve: easeOutQuart,
        ),
        FadeEffect(
          begin: 0,
          end: 1,
          duration: medium,
          curve: easeOutQuart,
        ),
      ];

  static List<Effect> get slideInFromLeft => [
        SlideEffect(
          begin: const Offset(-1.0, 0),
          end: Offset.zero,
          duration: medium,
          curve: easeOutQuart,
        ),
        FadeEffect(
          begin: 0,
          end: 1,
          duration: medium,
          curve: easeOutQuart,
        ),
      ];

  static List<Effect> get bounceIn => [
        ScaleEffect(
          begin: const Offset(0.3, 0.3),
          end: const Offset(1.0, 1.0),
          duration: slow,
          curve: bounceOut,
        ),
        FadeEffect(
          begin: 0,
          end: 1,
          duration: fast,
          curve: Curves.easeOut,
        ),
      ];

  static List<Effect> get shimmer => [
        ShimmerEffect(
          duration: const Duration(milliseconds: 1500),
          color: Colors.white.withOpacity(0.6),
          angle: 0,
        ),
      ];

  static List<Effect> get pulse => [
        ScaleEffect(
          begin: const Offset(1.0, 1.0),
          end: const Offset(1.05, 1.05),
          duration: const Duration(milliseconds: 1000),
          curve: Curves.easeInOut,
        ),
      ];

  // Page transitions
  static Widget slideTransition({
    required Widget child,
    required Animation<double> animation,
    SlideDirection direction = SlideDirection.right,
  }) {
    Offset begin;
    switch (direction) {
      case SlideDirection.up:
        begin = const Offset(0, 1);
        break;
      case SlideDirection.down:
        begin = const Offset(0, -1);
        break;
      case SlideDirection.left:
        begin = const Offset(1, 0);
        break;
      case SlideDirection.right:
        begin = const Offset(-1, 0);
        break;
    }

    return SlideTransition(
      position: Tween<Offset>(
        begin: begin,
        end: Offset.zero,
      ).animate(CurvedAnimation(
        parent: animation,
        curve: easeOutQuart,
      )),
      child: child,
    );
  }

  static Widget scaleTransition({
    required Widget child,
    required Animation<double> animation,
  }) {
    return ScaleTransition(
      scale: Tween<double>(
        begin: 0.8,
        end: 1.0,
      ).animate(CurvedAnimation(
        parent: animation,
        curve: easeOutQuart,
      )),
      child: FadeTransition(
        opacity: animation,
        child: child,
      ),
    );
  }

  static Widget rotationTransition({
    required Widget child,
    required Animation<double> animation,
  }) {
    return RotationTransition(
      turns: Tween<double>(
        begin: 0.0,
        end: 1.0,
      ).animate(CurvedAnimation(
        parent: animation,
        curve: easeInOutCubic,
      )),
      child: child,
    );
  }

  // Staggered animations
  static List<Widget> staggeredFadeInUp({
    required List<Widget> children,
    Duration delay = const Duration(milliseconds: 100),
  }) {
    return children.asMap().entries.map((entry) {
      int index = entry.key;
      Widget child = entry.value;

      return child
          .animate(delay: delay * index)
          .fadeIn(duration: medium, curve: easeOutQuart)
          .slideY(
            begin: 0.3,
            end: 0,
            duration: medium,
            curve: easeOutQuart,
          );
    }).toList();
  }

  // Loading animations
  static Widget loadingDots({Color? color}) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: List.generate(3, (index) {
        return Container(
          width: 8,
          height: 8,
          margin: const EdgeInsets.symmetric(horizontal: 2),
          decoration: BoxDecoration(
            color: color ?? Colors.blue,
            shape: BoxShape.circle,
          ),
        ).animate(onPlay: (controller) => controller.repeat()).scaleXY(
              begin: 0.5,
              end: 1.0,
              duration: const Duration(milliseconds: 600),
              delay: Duration(milliseconds: index * 200),
              curve: Curves.easeInOut,
            );
      }),
    );
  }

  static Widget loadingSpinner({Color? color, double size = 24}) {
    return SizedBox(
      width: size,
      height: size,
      child: CircularProgressIndicator(
        strokeWidth: 2,
        valueColor: AlwaysStoppedAnimation<Color>(
          color ?? Colors.blue,
        ),
      ),
    );
  }

  // Feedback animations
  static Widget buttonPress({required Widget child}) {
    return child
        .animate(
          target: 0,
          autoPlay: false,
        )
        .scaleXY(end: 0.95, duration: const Duration(milliseconds: 100));
  }

  static Widget heartbeat({required Widget child}) {
    return child.animate(onPlay: (controller) => controller.repeat()).scaleXY(
          begin: 1.0,
          end: 1.1,
          duration: const Duration(milliseconds: 800),
          curve: Curves.easeInOut,
        );
  }

  // Text animations
  static Widget typewriter({
    required String text,
    required TextStyle style,
    Duration duration = const Duration(milliseconds: 2000),
  }) {
    return AnimatedBuilder(
      animation: AlwaysStoppedAnimation(1.0),
      builder: (context, child) {
        return Text(
          text,
          style: style,
        ).animate().custom(
              duration: duration,
              builder: (context, value, child) {
                final int count = (value * text.length).round();
                return Text(
                  text.substring(0, count),
                  style: style,
                );
              },
            );
      },
    );
  }

  // Gesture animations
  static Widget swipeToDelete({
    required Widget child,
    required VoidCallback onDelete,
  }) {
    return Dismissible(
      key: UniqueKey(),
      direction: DismissDirection.endToStart,
      onDismissed: (_) => onDelete(),
      background: Container(
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 20),
        decoration: BoxDecoration(
          color: Colors.red,
          borderRadius: BorderRadius.circular(16),
        ),
        child: const Icon(
          Icons.delete,
          color: Colors.white,
        ),
      ),
      child: child,
    );
  }
}

enum SlideDirection { up, down, left, right }

// Custom route transitions
class CustomPageRoute<T> extends PageRouteBuilder<T> {
  final Widget child;
  final PageTransitionType transitionType;

  CustomPageRoute({
    required this.child,
    this.transitionType = PageTransitionType.slideRight,
    super.settings,
  }) : super(
          pageBuilder: (context, animation, _) => child,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            switch (transitionType) {
              case PageTransitionType.fade:
                return FadeTransition(opacity: animation, child: child);
              case PageTransitionType.scale:
                return AppAnimations.scaleTransition(
                  child: child,
                  animation: animation,
                );
              case PageTransitionType.slideRight:
                return AppAnimations.slideTransition(
                  child: child,
                  animation: animation,
                  direction: SlideDirection.right,
                );
              case PageTransitionType.slideLeft:
                return AppAnimations.slideTransition(
                  child: child,
                  animation: animation,
                  direction: SlideDirection.left,
                );
              case PageTransitionType.slideUp:
                return AppAnimations.slideTransition(
                  child: child,
                  animation: animation,
                  direction: SlideDirection.up,
                );
              case PageTransitionType.slideDown:
                return AppAnimations.slideTransition(
                  child: child,
                  animation: animation,
                  direction: SlideDirection.down,
                );
            }
          },
          transitionDuration: AppAnimations.medium,
        );
}

enum PageTransitionType {
  fade,
  scale,
  slideRight,
  slideLeft,
  slideUp,
  slideDown,
}

// Animation helper class instead of mixin
class AnimationHelper {
  static AnimationController createController({
    required TickerProvider vsync,
    Duration duration = const Duration(milliseconds: 300),
  }) {
    return AnimationController(duration: duration, vsync: vsync);
  }

  static Animation<double> createCurvedAnimation({
    required AnimationController controller,
    Curve curve = Curves.easeOutQuart,
  }) {
    return CurvedAnimation(parent: controller, curve: curve);
  }
}
