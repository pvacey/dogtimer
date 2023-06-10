<script>
  import { onMount } from 'svelte';
  let port = import.meta.env.DEV ? 3000: window.location.port;
  let base_url = `${window.location.protocol}//${window.location.hostname}:${port}`
  let previousLast = new Date().getTime();
  let lastTime = new Date().getTime();
  let message = 'loading...';
  let hours = '';
  let minutes = '';
  let diff = 0;
  let height = window.innerHeight/3;
  export let name;

  const updateClock = () => {
    setTimeout(() => {
      let now = new Date().getTime();
      diff = (now - lastTime)/1000;
      hours   = Math.floor(diff/60/60);
      minutes = Math.floor(diff/60%60);

      if (hours < 1) {
        message = `${minutes} mins`;
      }
      else if (hours === 1) {
        message = `${hours} hour<br>${minutes} mins`
      }
      else
        message = `${hours} hours<br>${minutes} mins`

      if (now % 3 === 0) {
        getLatestTime()
      }
      height = window.innerHeight/3;
      updateClock()
    }, 1000);
  }

  const undo = () => {
    message = ''
    lastTime = previousLast;
    setLatestTime();
  }

  const reset = () => {
    if (diff > 300) {
      message = "0 mins";
      previousLast = lastTime;
      lastTime = new Date().getTime();
      setLatestTime();
    }
  }

  async function getLatestTime() {
    const resp1 = await fetch(`${base_url}/time/${name}`)
      resp1.json().then(data => {
        lastTime = data.Time;
    });
    const resp2  = await fetch(`${base_url}/time/${name}_previous`)
      resp2.json().then(data => {
        previousLast = data.Time;
    });
  } 

  async function setLatestTime() {
    await fetch(`${base_url}/time/${name}/${lastTime}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    })
    await fetch(`${base_url}/time/${name}_previous/${previousLast}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    })
  } 
    
  onMount(() => {
    updateClock();
    getLatestTime();

  });
</script>

<button class="card {name}" style="height:{height*.78}px" on:click={reset}>
  {@html message}
</button>
<button class="undo" style="height:{height*.2}px" on:click={undo}>
  undo
</button>
