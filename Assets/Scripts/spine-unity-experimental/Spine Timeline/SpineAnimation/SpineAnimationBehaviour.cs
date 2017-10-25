using System;
using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.Timeline;
using Spine;
using Spine.Unity;
using System.Collections.Generic;

namespace Spine.Unity.Playables {

	using Animation = Spine.Animation;

	[Serializable]
	public class SpineAnimationBehaviour : PlayableBehaviour {
		public SkeletonDataAsset skeletonDataAsset;

		[SpineAnimation(dataField:"skeletonDataAsset")]
		public string animationName;
		public bool loop;

		[Range(0f, 1f)]
		public float eventThreshold, attachmentThreshold, drawOrderThreshold;

		internal Animation animation;
		//internal SpineAnimationBehaviour previous;

//		internal readonly ExposedList<int> timelineData = new ExposedList<int>();
//		internal readonly ExposedList<SpineAnimationBehaviour> timelineDipMix = new ExposedList<SpineAnimationBehaviour>();
//		internal readonly ExposedList<float> timelinesRotation = new ExposedList<float>();
//
//		internal bool HasTimeline (int id) {
//			var timelines = animation.timelines.Items;
//			for (int i = 0, n = animation.timelines.Count; i < n; i++)
//				if (timelines[i].PropertyId == id) return true;
//			return false;
//		}

		public void EnsureInitialize (SkeletonData data) {
			if (animation == null) {
				animation = data.FindAnimation(animationName);
				//this.previous = previous;
			}
		}

	}

}

