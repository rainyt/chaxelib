import zygame.utils.SceneManager;
import zygame.display.Quad;
import zygame.core.Start;

class Main extends Start {
	static function main() {
		Start.initApp(Main, 0, 0, false);
	}

	override function init() {
		super.init();
		// ...
		this.engine.backgroundColor = 0xffffff;
		SceneManager.replaceScene(MainScene);
	}
}
