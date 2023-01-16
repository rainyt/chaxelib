import zygame.display.Label;
import zygame.display.Quad;
import zygame.display.Scene;

class MainScene extends Scene {
	override function onInit() {
		super.onInit();
		var quad:Quad = new Quad(100, 100);
		quad.left = 0;
		quad.right = 0;
		this.addChild(quad);

		var label = new Label("CHaxelib v1.0.0\n中国国内镜像库下载", this);
		label.setSize(56);
		label.textAlign = Center;
		label.setColor(0x0);
		label.centerX = -label.textWidth / 2;
		label.centerY = 0;
	}
}
