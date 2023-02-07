<script lang="ts">
	import { invalidate, invalidateAll } from '$app/navigation';
	import { apiUrl, type APIResponse } from '$lib/api';
	import type { PageServerData } from './$types';
	import Output from './Output.svelte';
	import SaveImage from './SaveImage.svelte';
	import ViewDrive from './ViewDrive.svelte';

	export let data: PageServerData;

	let websiteUrl = '';
	let optionsHeight = 800;
	let optionsWidth = 1280;
	let optionsFullpage = false;
	let optionsImageType = 'png';
	let optionsDriver = 'playwright';
	let optionsWhitehole = '';
	let optionsSaveDrive = false;
	let optionsNoOutput = false;

	let showOptions = true;
	let show = false;
	let error = false;
	let errorData: APIResponse | null = null;
	let fetching = false;
	let done = false;

	let imageUrl: string = '';

	const fetchScreenshot = async () => {
		if (!websiteUrl) return;

		imageUrl = '';
		show = true;
		fetching = true;
		done = false;

		if (optionsDriver === 'playwright') {
			optionsImageType = 'png';
		}

		if (!optionsSaveDrive) {
			optionsNoOutput = false;
		}

		const params = new URLSearchParams({
			height: optionsHeight.toString(),
			width: optionsWidth.toString(),
			imageType: optionsImageType,
			fullPage: optionsFullpage.toString(),
			driver: optionsDriver,
			whiteholeUrl: optionsWhitehole,
			saveToDrive: optionsSaveDrive.toString(),
			saveNoOutput: optionsNoOutput.toString()
		});
		const paramQueries = params.toString();
		if (paramQueries == '') {
			return;
		}
		const url = paramQueries != '' ? `/screenshot?${paramQueries}` : '/screenshot';

		const r = await fetch(apiUrl + url, {
			method: 'POST',
			headers: {
				'content-type': 'application/json'
			},
			body: JSON.stringify({ website: websiteUrl })
		});

		fetching = false;

		if (!r.ok) {
			error = true;
			errorData = await r.json();
			return;
		}

		error = false;
		errorData = null;
		done = true;

		if (optionsSaveDrive) {
			// re-fetch if screenshot is saved to drive
			invalidateAll();
		}

		if (optionsNoOutput) {
			return;
		}

		const imageData = await r.blob();
		imageUrl = URL.createObjectURL(imageData);
	};
</script>

<header class="py-4 w-5/6 mx-auto text-center">
	<span class="inline-flex items-center">
		<img src="/favicon.png" alt="Koyo" class="object-fit h-12 w-12" />
		<h1 class="font-black text-xl text-blue-500">Koyo</h1>
	</span>
	<p class="text-neutral-600">Website screenshot service api</p>
	<a
		title="View Github Repo"
		href="https://github.com/tbdsux/koyo"
		target="_blank"
		rel="noopener noreferrer"
		class="inline-flex mx-auto mt-2 opacity-80 hover:opacity-100 duration-300"
	>
		<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-5 w-5">
			<path
				d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
			/>
		</svg>
	</a>
</header>

<hr class="w-1/2 mx-auto my-3 border-gray-100" />

<div>
	<div class="flex flex-col text-left w-full sm:w-11/12 md:w-5/6 lg:w-4/5 xl:w-3/5 mx-auto">
		<label for="website" class="text-neutral-700">Enter website url:</label>
		<div class="flex items-center flex-col md:flex-row">
			<input
				bind:value={websiteUrl}
				type="text"
				id="website"
				name="website"
				class="m-0.5 py-2 px-5 rounded-xl border-2 outline-none hover:border-blue-500 focus:border-blue-500 w-full"
				placeholder="https://www.deta.space"
			/>

			<button
				on:click={fetchScreenshot}
				disabled={fetching}
				class="m-0.5 inline-flex items-center text-sm bg-blue-400 hover:bg-blue-500 py-2 px-8 rounded-xl text-white duration-300 disabled:opacity-80 disabled:hover:bg-blue-400"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25"
					/>
				</svg>

				<span class="ml-1">
					{#if fetching}
						Fetching...
					{:else}
						Screenshot
					{/if}</span
				>
			</button>

			<ViewDrive bind:data />
		</div>

		<div class="mt-2">
			<button class="text-xs text-gray-600" on:click={() => (showOptions = !showOptions)}>
				{showOptions ? 'Hide' : 'Show'} options</button
			>

			{#if showOptions}
				<div class="px-2 py-4 rounded-xl border border-gray-100">
					<div class="flex items-center justify-center flex-wrap">
						<div class="inline-flex flex-col m-2">
							<label for="width" class="text-sm text-gray-700">Width</label>
							<input
								bind:value={optionsWidth}
								type="number"
								min={240}
								max={16834}
								name="width"
								id="width"
								class="py-2 px-3 text-sm border rounded-xl w-32 disabled:opacity-80"
							/>
						</div>

						<div class="inline-flex flex-col m-2">
							<label for="height" class="text-sm text-gray-700">Height</label>
							<input
								bind:value={optionsHeight}
								disabled={optionsFullpage}
								type="number"
								min={240}
								max={16834}
								name="height"
								id="height"
								class="py-2 px-3 text-sm border rounded-xl w-32"
							/>
						</div>

						<div class="inline-flex flex-col m-2">
							<label for="width" class="text-sm text-gray-700">Image Type</label>
							<select
								bind:value={optionsImageType}
								class="py-2 px-3 text-sm border rounded-xl w-32 disabled:opacity-80 bg-white"
							>
								<option value="png">png</option>
								<option value="jpeg">jpeg</option>
								<option disabled={optionsDriver === 'playwright'} value="webp">webp</option>
							</select>
						</div>

						<div class="inline-flex flex-col m-2">
							<label for="width" class="text-sm text-gray-700">Driver</label>
							<select
								bind:value={optionsDriver}
								class="py-2 px-3 text-sm border rounded-xl w-32 disabled:opacity-80 bg-white"
							>
								<option value="playwright">playwright</option>
								<option value="puppeteer">puppeteer</option>
							</select>
						</div>

						<div class="inline-flex items-center flex-row m-2">
							<input
								bind:checked={optionsFullpage}
								type="checkbox"
								name="fullpage"
								id="fullpage"
								class="h-4 w-4 text-sm border rounded-xl"
							/>
							<label for="fullpage" class="text-sm ml-1 text-gray-700">Full Page</label>
						</div>
					</div>

					<hr class="my-4 mx-8 border-gray-100" />

					<div class="flex flex-wrap items-end justify-center">
						<div class="inline-flex items-center flex-row m-2">
							<input
								bind:checked={optionsSaveDrive}
								type="checkbox"
								name="fullpage"
								id="fullpage"
								class="h-4 w-4 text-sm border rounded-xl"
							/>
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<span
								on:click={() => (optionsSaveDrive = !optionsSaveDrive)}
								class="text-sm ml-1 text-gray-700 cursor-pointer"
								title="Do not return the image screenshot output.">Save to Drive (auto)</span
							>
						</div>

						{#if optionsSaveDrive}
							<div class="inline-flex items-center flex-row m-2">
								<input
									bind:checked={optionsNoOutput}
									type="checkbox"
									name="fullpage"
									id="fullpage"
									class="h-4 w-4 text-sm border rounded-xl"
								/>
								<!-- svelte-ignore a11y-click-events-have-key-events -->
								<span
									on:click={() => (optionsNoOutput = !optionsNoOutput)}
									class="text-sm ml-1 text-gray-700 cursor-pointer"
									title="Do not return the image screenshot output.">No Output</span
								>
							</div>
						{/if}
					</div>

					<hr class="my-4 mx-8 border-gray-100" />

					<div class="flex items-center justify-center">
						<div class="flex flex-col w-11/12 md:w-5/6  2xl:w-1/2 mx-auto">
							<label for="whitehole" class="">
								<a
									href="https://alpha.deta.space/discovery/@mikhailsdv/black_hole-3kf"
									target="_blank"
									rel="noopener noreferrer"
									title="Add Whitehole integration"
									class="text-sm text-gray-700 inline-flex items-center hover:text-gray-900 hover:underline"
								>
									<span class="mr-1">Whitehole</span>

									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-4 w-4"
										viewBox="0 0 20 20"
										fill="currentColor"
									>
										<path
											d="M11 3a1 1 0 100 2h2.586l-6.293 6.293a1 1 0 101.414 1.414L15 6.414V9a1 1 0 102 0V4a1 1 0 00-1-1h-5z"
										/>
										<path
											d="M5 5a2 2 0 00-2 2v8a2 2 0 002 2h8a2 2 0 002-2v-3a1 1 0 10-2 0v3H5V7h3a1 1 0 000-2H5z"
										/>
									</svg>
								</a>
							</label>
							<input
								bind:value={optionsWhitehole}
								type="text"
								placeholder="Whitehole integration url"
								name="whitehole"
								id="whitehole"
								class="py-2 px-3 text-sm border rounded-xl w-full"
							/>
						</div>
					</div>
				</div>
			{/if}
		</div>

		<div class="my-3">
			<a
				class="text-sm underline text-gray-600 hover:text-gray-700"
				href="https://deta.space/docs/en/basics/micros#api-keys"
				target="_blank"
				rel="noopener noreferrer">Where to get API Keys?</a
			>
		</div>

		<pre
			class="mt-4 overflow-auto text-left text-sm bg-gray-100 text-gray-700 py-2 px-3 rounded-xl">
{`curl -X POST \\
	'${
		data.baseUrl
	}/screenshot?height=${optionsHeight}&width=${optionsWidth}&imageType=${optionsImageType}&fullPage=${optionsFullpage}&saveToDrive=${optionsSaveDrive}${
				optionsSaveDrive ? '&saveNoOutput=' + optionsNoOutput.toString() : ''
			}${optionsWhitehole ? '&whiteholeUrl=' + optionsWhitehole : ''}' \\
	--header 'Content-Type: application/json' \\
	--header 'X-Space-App-Key: your-space-app-api-key' \\
	--data-raw '{"website":"${websiteUrl}"}' \\
	--output screenshot.${optionsImageType}`}
		</pre>
	</div>

	{#if show}
		<div class="mt-12">
			{#if imageUrl != ''}
				<SaveImage bind:imageUrl bind:websiteUrl bind:optionsImageType />
			{/if}

			<div class="mt-2 h-screen w-full relative overflow-auto mb-20">
				{#if fetching}
					<div class="h-full w-full animate-pulse bg-gray-300" />
				{:else if error}
					<pre class="text-sm text-left p-4 bg-gray-50">{JSON.stringify(errorData, null, 4)}</pre>
				{:else if done}
					<Output bind:imageUrl bind:websiteUrl {optionsNoOutput} />
				{/if}
			</div>
		</div>
	{/if}
</div>
