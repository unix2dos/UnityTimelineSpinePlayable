using UnityEditor;
using UnityEngine;
using Spine;
using Spine.Unity;
using Spine.Unity.Playables;

[CustomPropertyDrawer(typeof(SpineAnimationBehaviour))]
public class SpineAnimationDrawer : PropertyDrawer {
	public override float GetPropertyHeight (SerializedProperty property, GUIContent label) {
		const int fieldCount = 3;
		return fieldCount * EditorGUIUtility.singleLineHeight;
	}

	public override void OnGUI (Rect position, SerializedProperty property, GUIContent label) {
		SerializedProperty skeletonDataAssetProp = property.FindPropertyRelative("skeletonDataAsset");
		SerializedProperty animationNameProp = property.FindPropertyRelative("animationName");
		SerializedProperty loopProp = property.FindPropertyRelative("loop");
		//SerializedProperty mixPoseProp = property.FindPropertyRelative("mixPose");

		Rect singleFieldRect = new Rect(position.x, position.y, position.width, EditorGUIUtility.singleLineHeight);
		EditorGUI.PropertyField(singleFieldRect, skeletonDataAssetProp);

		singleFieldRect.y += EditorGUIUtility.singleLineHeight;
		EditorGUI.PropertyField(singleFieldRect, animationNameProp);

		singleFieldRect.y += EditorGUIUtility.singleLineHeight;
		EditorGUI.PropertyField(singleFieldRect, loopProp);
	}
}
