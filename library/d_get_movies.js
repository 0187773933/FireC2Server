( async ()=> {
	// https://www.disneyplus.com/movies/9f7c38e5-41c3-47b4-b99e-b5b3d2eb95d4
	function sleep( ms ) { return new Promise( resolve => setTimeout( resolve , ms ) ); }
	async function scroll_to_bottom() {
		let lastScrollHeight = 0;
		while (true) {
			window.scrollTo(0, document.documentElement.scrollHeight);
			await sleep(1500);
			let currentScrollHeight = document.documentElement.scrollHeight;
			if (lastScrollHeight === currentScrollHeight) {
				await sleep(2000);
				window.scrollTo(0, document.documentElement.scrollHeight);
				currentScrollHeight = document.documentElement.scrollHeight;
				if (lastScrollHeight === currentScrollHeight) {
					console.log("Reached the true bottom!");
					break;
				}
			}
			lastScrollHeight = currentScrollHeight;
		}
	}
	await scroll_to_bottom();

	let result = {};
	let nodes = document.querySelectorAll( 'div.gv2-asset > a[data-gv2elementvalue]' );
	nodes.forEach(node => {
		let firstChild = node.firstElementChild;
		if (firstChild && firstChild.hasAttribute('aria-label')) {
			let id = node.getAttribute('data-gv2elementvalue');
			let name = firstChild.getAttribute('aria-label');
			result[id] = name;
		}
	});
	// console.log( JSON.stringify( result ) );
	let entries = Object.entries(result);
	entries.sort((a, b) => a[1].localeCompare(b[1]));
	let result_string = "";
	entries.forEach(([id, value]) => {
		result_string += `${value} === ${id}\n`;
	});
	console.log( result_string )
})();