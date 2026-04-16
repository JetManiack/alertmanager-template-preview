export const gotemplate = {
  name: "gotemplate",
  startState: () => ({
    inTag: false,
    inComment: false
  }),
  token: (stream, state) => {
    // 1. Handle template comments: {{/* ... */}}
    if (state.inComment) {
      if (stream.match('*/}}') || stream.match('*/-}}')) {
        state.inComment = false;
        return 'comment';
      }
      stream.next();
      return 'comment';
    }

    // 2. Inside a tag: {{ ... }}
    if (state.inTag) {
      if (stream.eatSpace()) return null;

      // Close tag (with optional whitespace control)
      if (stream.match('}}') || stream.match('-}}')) {
        state.inTag = false;
        return 'punctuation';
      }

      // Strings
      if (stream.match(/^"(?:[^"\\]|\\.)*"/)) return 'string';
      if (stream.match(/^'(?:[^'\\]|\\.)*'/)) return 'string';
      if (stream.match(/^`[^`]*`/)) return 'string';

      // Keywords
      if (stream.match(/^(?:if|else|range|with|template|define|end|block)\b/)) return 'keyword';

      // Atoms
      if (stream.match(/^(?:true|false|nil|iota)\b/)) return 'atom';

      // Numbers
      if (stream.match(/^[0-9]+(\.[0-9]+)?/)) return 'number';

      // Variables / Fields starting with dot (e.g., .Field, .Labels.name)
      if (stream.match(/^\.[a-zA-Z_][a-zA-Z0-9_]*/)) return 'variableName';

      // Identifiers / Functions (e.g., humanize, toUpper)
      if (stream.match(/^[a-zA-Z_][a-zA-Z0-9_]*/)) return 'variableName';

      // Operators and other punctuation
      if (stream.match(/^[+\-*&%=<>!|:]+/)) return 'operator';
      if (stream.match(/^[[\]{}(),;.]/)) return 'punctuation';

      stream.next();
      return null;
    } 
    
    // 3. Outside of tags (plain text)
    else {
      // Check for start of a comment
      if (stream.match('{{/*') || stream.match('{{-/*')) {
        state.inComment = true;
        return 'comment';
      }
      // Check for start of a tag (with optional whitespace control)
      if (stream.match('{{-') || stream.match('{{')) {
        state.inTag = true;
        return 'punctuation';
      }
      // Just plain text, skip ahead to next {{ or end of line
      if (!stream.skipTo('{{')) {
        stream.skipToEnd();
      }
      return null;
    }
  }
};
