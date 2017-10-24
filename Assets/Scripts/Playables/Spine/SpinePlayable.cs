using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.Playables;
using Spine.Unity;

//一个clip一个实例
public class SpinePlayable : PlayableBehaviour
{
    //static变量表明是一个timeline生命周期
    static Dictionary<int, int> _dict = new Dictionary<int, int> ();//记录该clip对应的index
    static List<SkeletonAnimation> _list = new List<SkeletonAnimation>();//记录了所有的spine
    static  int _index = -1;

    private GameObject _spineObject;
    private string _actionName;
    private bool _actionHold;


    public void Initialize (GameObject spineObject, string name, bool actionHold)
    {
        _spineObject = spineObject;
        _actionName = name;
        _actionHold = actionHold;
    }


    public override void OnBehaviourPlay (Playable playable, FrameData info)
    {
        if (_spineObject != null) {
            var skeletonAnimation = _spineObject.GetComponent<SkeletonAnimation> ();
            if (skeletonAnimation != null && skeletonAnimation.AnimationState != null) {
                skeletonAnimation.AnimationState.SetAnimation (getPlayTrackIndex (), _actionName, true);
                if (_list.Contains (skeletonAnimation) == false) {
                    _list.Add (skeletonAnimation);
                }
            }
        }
    }

    public override void OnBehaviourPause (Playable playable, FrameData info)
    {
        if (_spineObject != null && _actionHold == false) {
            var skeletonAnimation = _spineObject.GetComponent<SkeletonAnimation> ();
            if (skeletonAnimation != null && skeletonAnimation.AnimationState != null) {
                skeletonAnimation.AnimationState.SetEmptyAnimation (getOverTrackIndex (), 0.2f);
            }
        }
    }


    //相同的clip index一样 不同的clip index不一样 这样就可以同时播多个
    public int getPlayTrackIndex ()
    {
        int code = this.GetHashCode ();
        if (_dict.ContainsKey (code) == false) {
            ++_index;
            _dict.Add (code, _index);
//            Debug.Log (this.GetHashCode() +  " play action : "  + _actionName  +  " " + _index);
        }
        return _dict [code];
    }


    public int getOverTrackIndex ()
    {
        int code = this.GetHashCode ();
        int index = 0;
        if (_dict.ContainsKey (code)) {
            index = _dict [code];
            _dict.Remove (code);
//            Debug.Log (this.GetHashCode () + " stop action : " + _actionName + " " + index);
        } else {
            //此处状态说明新的timeline开始了, 是unity的整个timeline, 要清除以前的hold住的动画
            foreach (var l in _list) {
                l.AnimationState.SetEmptyAnimations (0.2f);
            }
            _list.Clear ();
            _dict.Clear ();
            _index = -1;
        }
        return index;
    }

}
