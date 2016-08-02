import {
	blue500, blue700, amber200,
	grey100, grey300, grey400, grey500,
	white, darkBlack, fullBlack,
} from 'material-ui/styles/colors';
import spacing from 'material-ui/styles/spacing';
import {fade} from 'material-ui/utils/colorManipulator';

export default {
	spacing: spacing,
	fontFamily: 'Roboto, sans-serif',
	palette: {
		primary1Color: blue500,
		primary2Color: blue700,
		primary3Color: grey400,
		accent1Color: amber200,
		accent2Color: grey100,
		accent3Color: grey500,
		textColor: darkBlack,
		alternateTextColor: white,
		canvasColor: white,
		borderColor: grey300,
		disabledColor: fade(darkBlack, 0.3),
		pickerHeaderColor: blue500,
		clockCircleColor: fade(darkBlack, 0.07),
		shadowColor: fullBlack,
	}
};