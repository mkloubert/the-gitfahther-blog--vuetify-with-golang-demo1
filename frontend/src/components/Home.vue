
  <script lang="ts">
  // MIT License
  //
  // Copyright (c) 2024 Marcel Joachim Kloubert (https://marcel.coffee)
  //
  // Permission is hereby granted, free of charge, to any person obtaining a copy
  // of this software and associated documentation files (the "Software"), to deal
  // in the Software without restriction, including without limitation the rights
  // to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  // copies of the Software, and to permit persons to whom the Software is
  // furnished to do so, subject to the following conditions:
  //
  // The above copyright notice and this permission notice shall be included in all
  // copies or substantial portions of the Software.
  //
  // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  // IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  // FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  // AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  // LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  // OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  // SOFTWARE.

  import axios from 'axios';

  export default {
    data: () => {
      return {
        isSendingText: false,
        lastStatus: '',
        textValue: '',
        voice: 'Alloy'
      };
    },
    
    methods: {
      async SendText() {        
        const textToSend = String(this.textValue).substring(0, 4000).trim();
        if (textToSend === "") {
          return;
        }

        const voice = String(this.voice).toLowerCase().trim() || "alloy";

        this.isSendingText = true;

        try {
          const {
            data,
            status
          } = await axios.post('http://localhost:4000/', {
            text: textToSend,
            voice
          }, {
            responseType: 'text'
          });

          if (status !== 200) {
            throw new Error(`Unexpected response: ${data}`);
          }

          const audio = new Audio(data);
          audio.volume = 0.5;
          audio.play();

          console.log('Playing audio');
        } catch (error) {
          console.error('[ERROR]', 'SendText()', error)
        } finally {
          this.isSendingText = false;
        }
      }
    }
  }
</script>

<template>
  <v-container class="fill-height">
    <v-responsive
      class="align-centerfill-height mx-auto"
      max-width="900"
    >
      <v-row>
        <v-col cols="12">
          <h1>Text-to-speech demo</h1>

          <v-textarea label="Your text" v-model="textValue" :maxlength="4000" counter></v-textarea>

          <v-select
            label="Voice"
            :items="['Alloy', 'Echo', 'Fable', 'Onyx', 'Nova', 'Shimmer']"
            v-model="voice"
          ></v-select>

          <v-btn @click="SendText" size="x-large" block :disabled="isSendingText || textValue.trim() === ''">
            Read out text
          </v-btn>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>
