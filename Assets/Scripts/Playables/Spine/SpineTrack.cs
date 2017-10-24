using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.Timeline;
using Spine.Unity;

[TrackColor (1f, 0f, 0f)]
[TrackClipType (typeof(SpineAsset))]
//[TrackBindingType (typeof(SkeletonAnimation))]
public class SpineTrack : TrackAsset
{
    public override Playable CreateTrackMixer (PlayableGraph graph, GameObject go, int inputCount)
    {
        return ScriptPlayable<SpinePlayable>.Create (graph, inputCount);
    }


//    public override void GatherProperties (PlayableDirector director, IPropertyCollector driver)
//    {
//        #if UNITY_EDITOR
//        SkeletonAnimation trackBinding = director.GetGenericBinding (this) as SkeletonAnimation;
//        if (trackBinding == null)
//            return;
//
//        var serializedObject = new UnityEditor.SerializedObject (trackBinding);
//        var iterator = serializedObject.GetIterator ();
//        while (iterator.NextVisible (true)) {
//            if (iterator.hasVisibleChildren)
//                continue;
//
//            driver.AddFromName<SkeletonAnimation> (trackBinding.gameObject, iterator.propertyPath);
//        }
//        #endif
//        base.GatherProperties (director, driver);
//    }
}
