using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.Timeline;
using System.Collections.Generic;

using Spine.Unity;

namespace Spine.Unity.Playables {
	
	[TrackColor(0.855f, 0.8623f, 0.87f)]
	[TrackClipType(typeof(SpineSkeletonFlipClip))]
	[TrackBindingType(typeof(SkeletonRenderer))]
	public class SpineSkeletonFlipTrack : TrackAsset {
		public override Playable CreateTrackMixer (PlayableGraph graph, GameObject go, int inputCount) {
			return ScriptPlayable<SpineSkeletonFlipMixerBehaviour>.Create(graph, inputCount);
		}

		public override void GatherProperties (PlayableDirector director, IPropertyCollector driver) {
#if UNITY_EDITOR
			SkeletonRenderer trackBinding = director.GetGenericBinding(this) as SkeletonRenderer;
			if (trackBinding == null)
				return;

			var serializedObject = new UnityEditor.SerializedObject(trackBinding);
			var iterator = serializedObject.GetIterator();
			while (iterator.NextVisible(true)) {
				if (iterator.hasVisibleChildren)
					continue;

				driver.AddFromName<SkeletonRenderer>(trackBinding.gameObject, iterator.propertyPath);
			}
#endif
			base.GatherProperties(director, driver);
		}
	}
}