using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.Playables;
using Spine.Unity;
using Spine;
using UnityEditor;
using UnityEngine.Timeline;


public class SpineAsset : PlayableAsset, ITimelineClipAsset
{
    public ExposedReference<GameObject> spineObject;

    [SpineAnimation(dataField:"_skeletonDataAsset")]
    public string actionName;

    [ContextMenuItem ("ActionDuration", "getActionDuration")]
    public bool actionKeep = false;//播放完毕是否保持


    //spine对象
    private SkeletonAnimation _spine = null;
    [HideInInspector]
    public SkeletonDataAsset _skeletonDataAsset;


//    [ContextMenuItem ("ActionNames", "getActionNames")]
//    public string actionName;
//    private void getActionNames ()
//    {
//        if (_spine != null &&  _spine.AnimationState != null) {
//            string strNames = "";
//            ExposedList<Spine.Animation> animations = _spine.AnimationState.Data.skeletonData.animations;
//            for (int i = 0, n = animations.Count; i < n; i++) {
//                strNames = strNames + animations.Items [i].name + "  ";
//            }
//            Debug.Log ("Actions:--->   " + strNames);
//        }
//    }

    private void getActionDuration ()
    {
        if (_spine != null &&  _spine.AnimationState != null) {
            var duration = _spine.AnimationState.Data.skeletonData.FindAnimation (actionName).Duration;
            Debug.Log (actionName +  "      Duration: " + duration);
        }
    }


    public override Playable CreatePlayable (PlayableGraph graph, GameObject owner)
    {
        var playable = ScriptPlayable<SpinePlayable>.Create (graph);
        var spine = spineObject.Resolve (graph.GetResolver ());
//        var abc = owner.GetComponent<PlayableDirector> ();
        if (spine != null)
        {
            _spine = spine.GetComponent<SkeletonAnimation> ();
            _skeletonDataAsset = _spine.SkeletonDataAsset;
        }
        playable.GetBehaviour ().Initialize (spine, actionName, actionKeep);
        return playable;
    }


    public ClipCaps clipCaps
    {
        get { return ClipCaps.None; }
    }
}

