<!--
// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
-->
<!--
pdf-text-layer: container for the retrievable text content of a pdf

props:
  sentenceHighlight: option to highlight the sentence that each match resides in
  page: the pdfjs page object to render text of
  textLayerDimStyle: an object representing the dimensions of the text
    container, usually containing the left and top of the corresponding canvas,
    and the height and width of the page's viewport.
  scale: the zoom level of the document

public methods:
  refresh(): request the component to rerender text, usually done when changing
    scale. This only marks the component to rerender on next update. This is
    partly done so expensive operations are not done during parent execution,
    and so we can ensure all dimensional data is updated before rendering.

events:
  'rendered': emitted upon completed rendering of the text
  'matches', { pageNum, matches }: emitted upon finishing the search for matches
-->
<template>
  <div :style="textLayerDimStyle" class="text-layer" :ref="textLayerID" />
</template>

<script>
import pdfjs from 'pdfjs-dist/webpack'
import { mapGetters } from 'vuex'

export default {
  name: 'pdf-text-layer',
  props: {
    sentenceHighlight: {
      type: Boolean,
      default: true
    },
    page: Object,
    textLayerDimStyle: Object,
    scale: Number
  },
  data () {
    return {
      textContent: null,
      renderTextLayerTask: null,
      textSpans: [],
      textContentItemsStr: [],
      joinedContent: '',
      joinedContentLower: '',
      matches: [],
      sentenceBounds: [],
      refreshTextLayer: false
    }
  },
  computed: {
    textLayerID () {
      return 'text-container-' + this.page.pageNumber
    },
    pixelRatio () {
      return window.devicePixelRatio || 1
    },
    sentenceStyle () {
      return this.sentenceHighlight ? 'sentenceOn' : 'sentenceOff'
    },
    ...mapGetters(['currentSearch'])
  },
  methods: {
    renderText () {
      if (this.textContent) {
        this.$refs[this.textLayerID].innerHTML = ''
        this.textSpans = []
        this.textContentItemsStr = []
        pdfjs.renderTextLayerTask = pdfjs.renderTextLayer({
          textContent: this.textContent,
          viewport: this.page.getViewport({
            scale: this.scale / this.pixelRatio
          }),
          container: this.$refs[this.textLayerID],
          textDivs: this.textSpans,
          textContentItemsStr: this.textContentItemsStr
        })
        this.$emit('rendered')
        this.refreshTextLayer = false
        this.matches = this.findMatches()
        this.highlightMatches()
        this.sendMatches()
      }
    },
    refresh () {
      this.refreshTextLayer = true
    },
    findMatches () {
      if (this.currentSearch.length === 0) return []
      const search = compileSearchTerms(this.currentSearch)
      this.textContentItemsStr = this.textContentItemsStr.map(str => {
        if (/\S/.test(str)) {
          return str
        } else {
          return ''
        }
      })
      this.joinedContent = this.textContentItemsStr.join('')
      this.joinedContentLower = this.joinedContent.toLowerCase()
      const { textSpans, joinedContentLower } = this
      this.sentenceBounds = findSentences()
      const sentenceBounds = this.sentenceBounds

      // take search query string and turn it into an array of terms, combining
      // terms with quotes.
      function compileSearchTerms (searchQuery) {
        const regex = /[^\s"]+|"([^"]*)"/g
        let result = []
        let match = regex.exec(searchQuery)
        while (match !== null) {
          result.push(match[1] ? match[1] : match[0])
          match = regex.exec(searchQuery)
        }
        result = result
          .filter(term => term.length > 0) // just in case
          .map(term => term.toLowerCase())
        return result
      }

      function getSpanFromJump (
        globalDelta,
        globalStart,
        spanStart,
        localStart
      ) {
        const globalEnd = globalStart + globalDelta
        while (globalDelta > 0 && spanStart < textSpans.length) {
          let toNextSpan = textSpans[spanStart].innerText.length - localStart
          if (toNextSpan < globalDelta) {
            spanStart++
            localStart = 0
            globalDelta -= toNextSpan
          } else {
            localStart += globalDelta
            globalDelta = 0
          }
        }
        if (spanStart >= textSpans.length) {
          spanStart = textSpans.length - 1
        }
        return {
          span: spanStart,
          offset: localStart,
          global: globalEnd
        }
      }

      function findSentences () {
        const punct = ['.', '!', '?']
        let result = []
        let localOffset = 0
        let globalOffset = 0
        let span = 0
        let minOffset = -1
        do {
          minOffset = -1
          punct.forEach(p => {
            const candidate = joinedContentLower.indexOf(p, globalOffset)
            if (
              candidate !== -1 &&
              (candidate <= minOffset || minOffset === -1)
            ) {
              minOffset = candidate
            }
          })
          if (minOffset > -1) {
            const next = getSpanFromJump(
              minOffset - globalOffset,
              globalOffset,
              span,
              localOffset
            )
            result.push({
              start: {
                span: span,
                offset: localOffset,
                global: globalOffset
              },
              end: {
                span: next.span,
                offset: next.offset,
                global: next.global
              }
            })
            const nextStart = getSpanFromJump(
              1,
              next.global,
              next.span,
              next.offset
            )
            globalOffset = nextStart.global
            localOffset = nextStart.offset
            span = nextStart.span
          }
        } while (minOffset > 0)
        const offsetToEnd = joinedContentLower.length - globalOffset
        if (offsetToEnd > 0) {
          const final = getSpanFromJump(
            joinedContentLower.length - globalOffset,
            globalOffset,
            span,
            localOffset
          )
          result.push({
            start: {
              span: span,
              offset: localOffset,
              global: globalOffset
            },
            end: {
              span: final.span,
              offset: final.offset,
              global: final.global
            }
          })
        }
        return result
      }

      function getNextMatch (globalIdx, localIdx, spanIdx) {
        // using the current offsets, get the next search term offset
        // try to keep this precedural without side effects
        let minOffset = -1
        let nextTermLength = -1
        search.forEach(currTerm => {
          if (currTerm.length === 0) return
          let candidateOffset = joinedContentLower.indexOf(currTerm, globalIdx)
          if (
            (candidateOffset <= minOffset && candidateOffset !== -1) ||
            (minOffset === -1 && candidateOffset !== -1)
          ) {
            // be sure to favor larger word e.g. for searching 'to' and 'tomorrow'
            if (
              candidateOffset === minOffset &&
              nextTermLength > currTerm.length
            ) {
              return
            }
            minOffset = candidateOffset // global offset
            nextTermLength = currTerm.length // keyword offset
          }
        })
        if (minOffset === -1) {
          // not found
          return null
        }
        // find sentence bounds and span idx
        // span idx
        let start = getSpanFromJump(
          minOffset - globalIdx,
          globalIdx,
          spanIdx,
          localIdx
        )
        // get end point
        let end = getSpanFromJump(
          nextTermLength,
          start.global,
          start.span,
          start.offset
        )

        // sentence bounds
        let sentenceIdx = 0
        let sentenceFound = false
        for (
          ;
          sentenceIdx < sentenceBounds.length && !sentenceFound;
          sentenceIdx++
        ) {
          const curr = sentenceBounds[sentenceIdx]
          sentenceFound =
            curr.start.global <= start.global && curr.end.global >= end.global
        }
        if (!sentenceFound) {
          // console.log('pdf-page: sentence not found: start:', start, 'end:', end)
        }
        return {
          start,
          end,
          sentence: sentenceIdx - 1
        }
      }

      let spanIdx = 0
      let localIdx = 0
      let globalIdx = 0
      let matches = []
      let match = getNextMatch(globalIdx, localIdx, spanIdx)
      while (match) {
        matches.push(match)
        spanIdx = match.end.span
        localIdx = match.end.offset
        globalIdx = match.end.global
        match = getNextMatch(globalIdx, localIdx, spanIdx)
      }
      return matches
    },
    sendMatches () {
      let matchContexts = []
      for (const matchIdx in this.matches) {
        let match = this.matches[matchIdx]
        let sentence = this.sentenceBounds[match.sentence]
        let sentenceText = this.joinedContent.substring(
          sentence.start.global,
          sentence.end.global
        )
        let span = this.textSpans[this.matches[matchIdx].start.span]
        matchContexts.push({
          sentence: {
            text: sentenceText,
            start: sentence.start,
            end: sentence.end
          },
          match,
          span
        })
      }
      this.$emit('matches', {
        pageNum: this.page.pageIndex,
        matches: matchContexts
      })
    },
    highlightMatches () {
      const sentSet = this.matches.reduce((acc, curr) => {
        return acc.add(curr.sentence)
      }, new Set())

      // create sentence segments
      const sentences = [...sentSet].map(idx => {
        return this.sentenceBounds[idx]
      })
      let segments = {}
      sentences.forEach(s => {
        let from = s.start.offset
        let to
        for (let spanIdx = s.start.span; spanIdx <= s.end.span; spanIdx++) {
          if (spanIdx === s.end.span) {
            to = s.end.offset
          } else {
            to = this.textContentItemsStr[spanIdx].length
          }
          if (!segments[spanIdx]) {
            segments[spanIdx] = []
          }
          segments[spanIdx].push({
            start: from,
            end: to,
            type: this.sentenceStyle
            // the 'type' portion of this object will determine the CSS style
            // used in appendTextChild()
          })
          from = 0
        }
      })
      // create keyword segments
      this.matches.forEach(match => {
        let segList = segments[match.start.span]
        let newStartSegList = []
        let i = 0
        while (i < segList.length && match.start.offset > segList[i].end) {
          newStartSegList.push(segList[i])
          i++
        }
        if (i === segList.length) {
          // console.log('keyword not found in segment')
          return
        }
        // cut sentence segment
        let from = segList[i].start
        let to = match.start.offset
        if (from < to) {
          newStartSegList.push({
            start: from,
            end: to,
            type: segList[i].type
          })
        }
        // start keyword segment
        from = match.start.offset
        if (match.start.span === match.end.span) {
          to = match.end.offset
        } else {
          to = this.textContentItemsStr[match.start.span].length
        }
        newStartSegList.push({
          start: from,
          end: to,
          type: 'keyword'
        })
        // push segments until we reach the end of the match
        for (let i = match.start.span + 1; i < match.end.span; i++) {
          const interSeg = [
            {
              start: 0,
              end: this.textContentItemsStr[i].length,
              type: 'keyword'
            }
          ]
          segments[i] = interSeg
        }
        // push tail end of keyword if we crossed spans
        let newTailSeg = []
        if (match.start.span !== match.end.span) {
          newTailSeg.push({
            start: 0,
            end: match.end.offset,
            type: 'keyword'
          })
        }
        // push tail end of sentence
        // find segment in match.end with segment end > match end offset
        segList = segments[match.end.span]
        i = 0
        while (i < segList.length && match.end.offset > segList[i].end) {
          i++
        }
        if (i === segList.length) {
          // console.log('match end not found in segment')
          return
        }
        from = match.end.offset
        to = segList[i].end
        newTailSeg.push({
          start: from,
          end: to,
          type: segList[i].type
        })
        // push rest of the last segment
        i++
        while (i < segList.length) {
          newTailSeg.push(segList[i])
          i++
        }
        // apply changes
        if (match.start.span === match.end.span) {
          newTailSeg.forEach(seg => {
            newStartSegList.push(seg)
          })
          segments[match.start.span] = newStartSegList
        } else {
          segments[match.start.span] = newStartSegList
          segments[match.end.span] = newTailSeg
        }
      })
      // fill in segments that have no highlighting
      for (let spanIdx in segments) {
        let spanContent = this.textContentItemsStr[spanIdx]
        let segmentList = segments[spanIdx]
        let newList = []
        let from = 0
        let to
        segmentList.forEach(seg => {
          to = seg.start
          if (from !== to) {
            newList.push({
              start: from,
              end: to,
              type: ''
            })
          }
          newList.push(seg)
          from = seg.end
        })
        const lastElementEnd = newList.slice(-1)[0].end
        if (lastElementEnd !== spanContent.length) {
          newList.push({
            start: lastElementEnd,
            end: spanContent.length,
            type: ''
          })
        }
        segments[spanIdx] = newList
      }

      // segments all documented. can now edit spans solely based on the object
      for (let spanIdx in segments) {
        let span = this.textSpans[spanIdx]
        span.textContent = ''
        segments[spanIdx].forEach(seg => {
          this.appendTextChild(spanIdx, seg.start, seg.end, seg.type)
        })
      }
    },
    appendTextChild (spanIdx, from, to, className, link) {
      let parentSpan = this.textSpans[spanIdx]
      let parentContent = this.textContentItemsStr[spanIdx]
      let substring = parentContent.substring(from, to)
      let textNode = document.createTextNode(substring)
      let childNode = document.createElement('span')
      if (className && !link) {
        childNode.className = className
        childNode.appendChild(textNode)
        parentSpan.appendChild(childNode)
      } else if (link) {
        childNode.className = className
        let linkNode = document.createElement('a')
        linkNode.setAttribute('href', link)
        linkNode.appendChild(textNode)
        childNode.appendChild(linkNode)
        parentSpan.appendChild(childNode)
      } else {
        parentSpan.appendChild(textNode)
      }
    }
  },
  mounted () {
    this.page.getTextContent().then(content => {
      this.textContent = content
      this.$nextTick(() => {
        this.renderText()
      })
    })
  },
  watch: {
    sentenceHighlight: 'renderText',
    refreshTextLayer (val) {
      if (val) {
        this.renderText()
      }
    }
  }
}
</script>

<style>
.text-layer {
  position: absolute;
  display: inline-block;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  opacity: 0.5;
  line-height: 1;
}
</style>

<style lang="scss">
.text-layer > span {
  color: transparent;
  position: absolute;
  white-space: pre;
  cursor: text;
  -webkit-transform-origin: 0% 0%;
  -moz-transform-origin: 0% 0%;
  -o-transform-origin: 0% 0%;
  -ms-transform-origin: 0% 0%;
  transform-origin: 0% 0%;
}

.text-layer ::selection {
  background: rgb(0, 0, 255);
}

.keyword {
  background-color: goldenrod;
}

.sentenceOn {
  background-color: $app-clr2;
}

.sentenceOff {
  background-color: transparent;
}

.link {
  background-color: green;
  cursor: pointer;
}
</style>
