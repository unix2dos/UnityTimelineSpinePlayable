using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.Playables;
using UnityEngine.UI;
using UnityEngine.Timeline;
//using TMPro;

public class ConversationPlayable : BasicPlayableBehaviour
{
//    [Header("Conversation Canvas")]
//    public ExposedReference<GameObject> canvasObject;
//    private GameObject _canvasObject;
//
//    [Header("Dialogue Speech")]
//    public ExposedReference<TextMeshProUGUI> dialogueTextDisplay;
//    private TextMeshProUGUI _dialogueTextDisplay;
//    public TMP_FontAsset fontAsset;
//    [Multiline(3)]
//    public string dialogueString;
//
//    [Header("Dialogue Box")]
//    public ExposedReference<Image> dialogueBoxDisplay;
//    private Image _dialogueBoxDisplay;
//    public Color dialogueBoxColor;

    public override void OnGraphStart(Playable playable){
//        _canvasObject = canvasObject.Resolve (playable.GetGraph ().GetResolver ());
//        _dialogueTextDisplay = dialogueTextDisplay.Resolve (playable.GetGraph ().GetResolver ());
//        _dialogueBoxDisplay = dialogueBoxDisplay.Resolve (playable.GetGraph ().GetResolver ());
    }

    public override void OnBehaviourPlay(Playable playable, FrameData info){
//        _canvasObject.SetActive (true);
//        _dialogueTextDisplay.font = fontAsset;
//        _dialogueTextDisplay.text = dialogueString;
//        _dialogueBoxDisplay.color = dialogueBoxColor;
    }

    public override void OnBehaviourPause(Playable playable, FrameData info){
//        _canvasObject.SetActive (false);
    }

}
