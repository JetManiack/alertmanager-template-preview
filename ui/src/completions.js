const alertmanagerFuncs = [
  "toUpper", "toLower", "title", "trimSpace", "join", "match", "safeHtml", 
  "safeUrl", "urlUnescape", "reReplaceAll", "stringSlice", "date", "tz", 
  "since", "humanizeDuration", "toJson", "list", "append", "dict"
].map(name => ({ label: name, type: "function" }));

const alertmanagerFields = [
  "Receiver", "Status", "Alerts", "NotificationReason", 
  "GroupLabels", "CommonLabels", "CommonAnnotations", "ExternalURL"
].map(name => ({ label: name, type: "variable" }));

const alertFields = [
  "Status", "Labels", "Annotations", "StartsAt", "EndsAt", 
  "GeneratorURL", "Fingerprint"
].map(name => ({ label: name, type: "variable" }));

const alertsMethods = [
  "Firing", "Resolved"
].map(name => ({ label: name, type: "method" }));

const kvMethods = [
  "SortedPairs", "Names", "Values"
].map(name => ({ label: name, type: "method" }));

export function createTemplateCompletionSource(alertData) {
  return (context) => {
    // Match a word before the cursor
    let word = context.matchBefore(/\w*/);
    
    if (!word || (word.from === word.to && !context.explicit)) {
      // Check if we just typed a dot
      const isJustDot = context.state.sliceDoc(context.pos - 1, context.pos) === ".";
      if (!isJustDot && !context.explicit) return null;
      
      if (isJustDot) {
        word = { from: context.pos, to: context.pos, text: "" };
      }
    }

    // Check if there is a dot before the word
    const isDot = word.from > 0 && context.state.sliceDoc(word.from - 1, word.from) === ".";
    
    if (isDot) {
      // Possible nested path: .Field.SubField
      const beforeDot = context.state.sliceDoc(0, word.from - 1);
      const parentMatch = beforeDot.match(/\.(\w+)$/);
      
      if (parentMatch) {
        const field = parentMatch[1];
        let source = null;
        if (field === "CommonLabels" && alertData?.commonLabels) source = alertData.commonLabels;
        else if (field === "GroupLabels" && alertData?.groupLabels) source = alertData.groupLabels;
        else if (field === "CommonAnnotations" && alertData?.commonAnnotations) source = alertData.commonAnnotations;
        
        if (source) {
          const keys = Object.keys(source).map(key => ({
            label: key,
            type: "property",
          }));
          return {
            from: word.from,
            to: context.pos,
            options: keys,
            validFor: /^\w*$/
          };
        }
      }

      // Default top-level fields
      return {
        from: word.from,
        to: context.pos,
        options: [...alertmanagerFields, ...alertFields, ...alertsMethods, ...kvMethods],
        validFor: /^\w*$/
      };
    }

    // No dot -> Suggest functions
    return {
      from: word.from,
      to: context.pos,
      options: alertmanagerFuncs,
      validFor: /^\w*$/
    };
  };
}
