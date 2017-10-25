using System;
using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.Timeline;
using Spine.Unity;
using Spine;
using System.Collections.Generic;

namespace Spine.Unity.Playables {
	public class SpineAnimationMixerBehaviour : PlayableBehaviour {
		SkeletonAnimation trackBindingSkeletonAnimation;
		Skeleton trackBindingSkeleton;


		static readonly Animation EmptyAnimation = new Animation("<empty>", new ExposedList<Timeline>(), 0);
		internal const int SUBSEQUENT = 0, FIRST = 1, DIP = 2, DIP_MIX = 3;

		readonly HashSet<int> propertyIDs = new HashSet<int>();
		readonly ExposedList<SpineAnimationBehaviour> mixingTo = new ExposedList<SpineAnimationBehaviour>();

		// NOTE: This function is called at runtime and edit time. Keep that in mind when setting the values of properties.
		public override void ProcessFrame (Playable playable, FrameData info, object playerData) {
			trackBindingSkeletonAnimation = playerData as SkeletonAnimation;
			if (!trackBindingSkeletonAnimation)
				return;

			int inputCount = playable.GetInputCount();

			if (trackBindingSkeleton == null)
				trackBindingSkeleton = trackBindingSkeletonAnimation.Skeleton;

			//trackBindingSkeleton.SetToSetupPose();
			for (int i = 0; i < inputCount; i++) {
				var inputPlayable = (ScriptPlayable<SpineAnimationBehaviour>)playable.GetInput(i); // The clip
				var clipBehaviourData = inputPlayable.GetBehaviour(); // the stateless data

				clipBehaviourData.EnsureInitialize(trackBindingSkeleton.Data);
				var animation = clipBehaviourData.animation;
				if (animation != null) animation.SetKeyedItemsToSetupPose(trackBindingSkeleton);
			}

			float totalWeight = 0f;
			int currentInputs = 0;

			for (int i = 0; i < inputCount; i++) {
				float inputWeight = playable.GetInputWeight(i);
				var inputPlayable = (ScriptPlayable<SpineAnimationBehaviour>)playable.GetInput(i); // The clip
				var clipBehaviourData = inputPlayable.GetBehaviour(); // the stateless data

				float time = (float)inputPlayable.GetTime(); // clip time.

				totalWeight += inputWeight;
//				clipBehaviourData.EnsureInitialize(trackBindingSkeleton.Data);
				var animation = clipBehaviourData.animation;
				if (animation != null) {

					if (!Mathf.Approximately(inputWeight, 0f)) {
						MixPose mixPose = currentInputs == 0 ? MixPose.Setup : MixPose.Current;
						MixDirection mixDirection = MixDirection.In;

						if (inputWeight < 1 && currentInputs > 1) {
							mixDirection = MixDirection.Out;
						}

						animation.Apply(trackBindingSkeleton, 0f, time, clipBehaviourData.loop, null, inputWeight, mixPose, mixDirection);
						Debug.LogFormat("Applying {0} at {1} as input [{2}] using {3} {4}", animation.Name, inputWeight, i, mixPose, mixDirection);

						currentInputs++;
					} else {
						//animation.Apply(trackBindingSkeleton, 0f, time, clipBehaviourData.loop, null, 0, MixPose.Current, MixDirection.Out);
						continue;
					}
				


				}
				// SPINETODO: Translate AnimationState into MixerBehaviour for robustness.
			}
		}

//		public override void OnGraphStart (Playable playable) {
//			Debug.Log("OnGraphStart");
//		}
//
//		public override void OnGraphStop (Playable playable) {
//			Debug.Log("OnGraphStop");
//
//			ScriptPlayable<SpineAnimationBehaviour> inputPlayable = (ScriptPlayable<SpineAnimationBehaviour>)playable.GetInput(0);
//			SpineAnimationBehaviour clipData = inputPlayable.GetBehaviour();
//			if (clipData == null)
//				return;
//
//			trackBindingSkeleton.SetToSetupPose(); // DEBUG
//			//trackBindingSkeletonAnimation.Update(0);
//		}
	}

}
