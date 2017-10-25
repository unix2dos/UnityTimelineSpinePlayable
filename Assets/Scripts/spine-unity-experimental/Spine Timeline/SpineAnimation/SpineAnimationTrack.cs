using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.Timeline;
using Spine.Unity;

namespace Spine.Unity.Playables {
	[TrackColor(0.9960785f, 0.2509804f, 0.003921569f)]
	[TrackClipType(typeof(SpineAnimationClip))]
	[TrackBindingType(typeof(SkeletonAnimation))]
	public class SpineAnimationTrack : TrackAsset {
		public override Playable CreateTrackMixer (PlayableGraph graph, GameObject go, int inputCount) {
			return ScriptPlayable<SpineAnimationMixerBehaviour>.Create(graph, inputCount);
		}

		public override void GatherProperties (PlayableDirector director, IPropertyCollector driver) {

			#if UNITY_EDITOR
			SkeletonAnimation trackBinding = director.GetGenericBinding(this) as SkeletonAnimation;
			if (trackBinding == null)
				return;

			var serializedObject = new UnityEditor.SerializedObject(trackBinding);
			var iterator = serializedObject.GetIterator();
			while (iterator.NextVisible(true)) {
				if (iterator.hasVisibleChildren)
					continue;

				driver.AddFromName<SkeletonAnimation>(trackBinding.gameObject, iterator.propertyPath);
			}
			#endif

			base.GatherProperties(director, driver);
		}
	}
}
