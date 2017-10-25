using System;
using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.Timeline;

namespace Spine.Unity.Playables {
	[Serializable]
	public class SpineAnimationClip : PlayableAsset, ITimelineClipAsset {
		public SpineAnimationBehaviour template = new SpineAnimationBehaviour();

		public ClipCaps clipCaps {
			get { return ClipCaps.Looping | ClipCaps.ClipIn | ClipCaps.SpeedMultiplier | ClipCaps.Blending; }
		}

		public override Playable CreatePlayable (PlayableGraph graph, GameObject owner) {
			var playable = ScriptPlayable<SpineAnimationBehaviour>.Create(graph, template);
			SpineAnimationBehaviour clone = playable.GetBehaviour();
			return playable;
		}
	}
}
